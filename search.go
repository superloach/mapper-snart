package mapper

import (
	"fmt"
	"strings"

	dg "github.com/bwmarrin/discordgo"

	"github.com/go-snart/snart/bot"
	"github.com/go-snart/snart/db"
	"github.com/go-snart/snart/db/admin"
	"github.com/go-snart/snart/route"
)

// Searcher is a wrapper that creates a function suitable for route.
func Searcher(b *bot.Bot, filt func(*Location) bool, limit int, msg string) func(ctx *route.Ctx) error {
	return func(ctx *route.Ctx) error {
		return Search(b.DB, ctx, admin.IsAdmin(b.DB)(ctx), filt, limit, msg)
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

	locs, err := GetLocations(d, ctx)

	for _, query := range queries {
		err = ctx.Session.ChannelTyping(ctx.Message.ChannelID)
		if err != nil {
			err = fmt.Errorf("typing %#v: %w", ctx.Message.ChannelID, err)
			Log.Warn(_f, err)
		}

		err = searchQuery(d, ctx, query, locs, limit, debug, msg)
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
		Log.Debug(_f, ps)

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

	if ps.Notes != "" {
		addField(embed, "Notes", ps.Notes)
	}

	if ps.Aliases != nil && len(ps.Aliases) > 0 {
		addField(embed, "Aliases", fmt.Sprintf("`%s`",
			strings.Join(ps.Aliases, "`, `")))
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
