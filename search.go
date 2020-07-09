package mapper

import (
	"fmt"
	"strconv"
	"strings"

	dg "github.com/bwmarrin/discordgo"
	"github.com/go-snart/snart/bot"
	"github.com/go-snart/snart/db"
	"github.com/go-snart/snart/route"
)

// Searcher is a wrapper that creates a function suitable for route.
func Searcher(b *bot.Bot, typ byte) func(ctx *route.Ctx) error {
	return func(ctx *route.Ctx) error {
		return Search(b.DB, ctx, b.Admin(ctx), typ)
	}
}

func cleanQueries(qs []string) []string {
	cqs := make([]string, 0, len(qs))

	for _, q := range qs {
		q = strings.TrimSpace(q)
		if len(q) == 0 {
			continue
		}

		cqs = append(cqs, q)
	}

	return cqs
}

func getPOIs(d *db.DB, m *dg.Message, typ byte) []*POI {
	_f := "getPOIs"

	d.Cache.Lock()

	poiCache := d.Cache.Get("mapper.poi").(db.Cache)

	bound := poiCache

	if d.Cache.Has("mapper.bound." + m.ChannelID) {
		bound = d.Cache.Get("mapper.bound." + m.ChannelID).(db.Cache)
	} else if d.Cache.Has("mapper.bound." + m.GuildID) {
		bound = d.Cache.Get("mapper.bound." + m.GuildID).(db.Cache)
	}

	d.Cache.Unlock()

	pois := []*POI{}

	bound.Lock()
	keys := bound.Keys()
	bound.Unlock()

	poiCache.Lock()
	for _, k := range keys {
		poi := poiCache.Get(k).(*POI)

		if typ == 0 || (len(poi.Pkmn) > 0 && poi.Pkmn[0] == typ) {
			pois = append(pois, poi)
		}
	}
	poiCache.Unlock()

	Log.Debugf(_f, "read %d", len(pois))

	return pois
}

// Search performs fuzzy searching on POIs.
func Search(
	d *db.DB, ctx *route.Ctx,
	admin bool, typ byte,
) error {
	_f := "Search"

	err := ctx.Session.ChannelTyping(ctx.Message.ChannelID)
	if err != nil {
		err = fmt.Errorf("typing %#v: %w", ctx.Message.ChannelID, err)
		Log.Warn(_f, err)
	}

	debug := false

	if admin {
		ctx.Flags.BoolVar(&debug, "debug", false, "print extra info")
	}

	_ = ctx.Flags.Parse()

	args := ctx.Flags.Args()
	queries := cleanQueries(strings.Split(strings.Join(args, " "), "+"))

	if len(queries) == 0 {
		rep1 := ctx.Reply()
		rep1.Content = "please specify a query.\nex: `" +
			ctx.CleanPrefix + ctx.Route.Name + " name of poi`"

		return rep1.Send()
	}

	msg := "POIs"
	limit := 100

	switch typ {
	case 'g':
		msg, limit = "Gyms", 25
	case 'p', 's':
		msg, limit = "PokéStops", 50
	}

	pois := getPOIs(d, ctx.Message, typ)

	for _, query := range queries {
		err = ctx.Session.ChannelTyping(ctx.Message.ChannelID)
		if err != nil {
			err = fmt.Errorf("typing %#v: %w", ctx.Message.ChannelID, err)
			Log.Warn(_f, err)
		}

		err = searchQuery(d, ctx, query, pois, limit, debug, msg)
		if err != nil {
			Log.Warn(_f, err)
		}
	}

	return nil
}

func addField(e *dg.MessageEmbed, name, value string) {
	e.Fields = append(e.Fields,
		&dg.MessageEmbedField{
			Name:  name,
			Value: value,
		},
	)
}

func searchQuery(
	d *db.DB, ctx *route.Ctx,
	query string, pois []*POI, limit int,
	debug bool, msg string,
) error {
	_f := "searchQuery"

	pss := search(clean(query), pois, 50, limit)
	if len(pss) == 0 {
		rep2 := ctx.Reply()
		rep2.Content = "no results found"

		return rep2.Send()
	}

	pg := NewWidget(ctx.Session, ctx.Message.ChannelID, ctx.Message.Author.ID)

	for i, ps := range pss {
		Log.Debugf(_f, "%#v\n", ps)

		ps.GetNeigh(d)

		e := mkEmbed(ps, i, len(pss), msg, nick(ctx.Message))

		if debug {
			debugEmbed(e, ps)
		}

		pg.Add(e)
	}

	go pg.Spawn()

	return nil
}

func mkEmbed(ps *poiScore, i, pssl int, msg, n string) *dg.MessageEmbed {
	embed := &dg.MessageEmbed{
		Title:       ps.Name,
		URL:         ps.URL(),
		Description: ps.Notes,
		Thumbnail: &dg.MessageEmbedThumbnail{
			URL:    ps.Image,
			Height: 100,
		},
		Footer: &dg.MessageEmbedFooter{
			Text: fmt.Sprintf(
				"%d/%d %s • %s",
				i+1, pssl, msg,
				n,
			),
		},
	}

	if ps.Neigh != nil {
		embed.Footer.Text = *ps.Neigh + " • " + embed.Footer.Text
	}

	if ps.Alias != nil && len(ps.Alias) > 0 {
		addField(embed,
			"Aliases",
			"`"+strings.Join(ps.Alias, "`, `")+"`",
		)
	}

	return embed
}

func debugEmbed(e *dg.MessageEmbed, ps *poiScore) {
	addField(e,
		"Flags",
		fmt.Sprintf("%q", map[string]string{
			"ingr": ps.Ingr,
			"pkmn": ps.Pkmn,
			"wzrd": ps.Wzrd,
		}),
	)

	addField(e,
		"Score",
		strconv.Itoa(ps.Score)+"%",
	)

	addField(e,
		"Link",
		ps.URL(),
	)
}
