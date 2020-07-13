package mapper

import (
	"fmt"
	"sort"
	"strings"

	dg "github.com/bwmarrin/discordgo"
	"github.com/go-snart/snart/bot"
	"github.com/go-snart/snart/db"
	"github.com/go-snart/snart/route"
)

// Searcher is a wrapper that creates a function suitable for route.
func Searcher(b *bot.Bot, filt func(*Location) bool, limit int, msg string) func(ctx *route.Ctx) error {
	return func(ctx *route.Ctx) error {
		return Search(b.DB, ctx, b.DB.Admin(ctx), filt, limit, msg)
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

func getLocations(d *db.DB, m *dg.Message, filt func(*Location) bool) []*Location {
	_f := "getLocations"

	d.Cache.Lock()

	locationCache := d.Cache.Get("mapper.location").(db.Cache)

	bound := locationCache

	if d.Cache.Has("mapper.bound." + m.ChannelID) {
		bound = d.Cache.Get("mapper.bound." + m.ChannelID).(db.Cache)
	} else if d.Cache.Has("mapper.bound." + m.GuildID) {
		bound = d.Cache.Get("mapper.bound." + m.GuildID).(db.Cache)
	}

	d.Cache.Unlock()

	locations := []*Location{}

	bound.Lock()
	keys := bound.Keys()
	bound.Unlock()

	locationCache.Lock()
	for _, k := range keys {
		location := locationCache.Get(k).(*Location)

		if filt(location) {
			locations = append(locations, location)
		}
	}
	locationCache.Unlock()

	Log.Debugf(_f, "read %d", len(locations))

	return locations
}

// Search performs fuzzy searching on Locations.
func Search(
	d *db.DB, ctx *route.Ctx, admin bool, filt func(*Location) bool, limit int, msg string,
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
		rep := ctx.Reply()
		rep.Content = fmt.Sprintf(
			"please specify a query.\nex: `%s%s [name of location]`",
			ctx.CleanPrefix, ctx.Route.Name)

		return rep.Send()
	}

	locations := getLocations(d, ctx.Message, filt)

	for _, query := range queries {
		err = ctx.Session.ChannelTyping(ctx.Message.ChannelID)
		if err != nil {
			err = fmt.Errorf("typing %#v: %w", ctx.Message.ChannelID, err)
			Log.Warn(_f, err)
		}

		err = searchQuery(d, ctx, query, locations, limit, debug, msg)
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
	query string, locations []*Location, limit int,
	debug bool, msg string,
) error {
	_f := "searchQuery"

	pss := search(clean(query), locations, 50, limit)
	if len(pss) == 0 {
		rep2 := ctx.Reply()
		rep2.Content = fmt.Sprintf("no %s found", msg)

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

func mkEmbed(ps *locationScore, i, pssl int, msg, n string) *dg.MessageEmbed {
	embed := &dg.MessageEmbed{
		Title: ps.Name,
		URL:   ps.URL(),
		Thumbnail: &dg.MessageEmbedThumbnail{
			URL:    ps.Image,
			Height: 100,
		},
		Footer: &dg.MessageEmbedFooter{
			Text: fmt.Sprintf(
				"%d/%d %s • %s",
				i+1, pssl, msg, n),
		},
	}

	if ps.Notes != nil {
		ks := []string(nil)
		for k := range ps.Notes {
			ks = append(ks, k)
		}

		sort.Strings(ks)

		for _, k := range ks {
			addField(embed, k, ps.Notes[k])
		}
	}

	if ps.Neigh != nil {
		embed.Footer.Text = *ps.Neigh + " • " + embed.Footer.Text
	}

	if ps.Alias != nil && len(ps.Alias) > 0 {
		addField(embed, "Aliases",
			"`"+strings.Join(ps.Alias, "`, `")+"`")
	}

	return embed
}

func debugEmbed(e *dg.MessageEmbed, ps *locationScore) {
	addField(e, "ID", ps.ID)

	addField(e, "Flags", fmt.Sprintf(
		"ingress: %s\npokemon: %s\nwizards: %s",
		ps.IngrType, ps.PkmnType, ps.WzrdType,
	))

	addField(e, "Score", fmt.Sprintf("%d%%", ps.Score))

	addField(e, "Link", ps.URL())
}
