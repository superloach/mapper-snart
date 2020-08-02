package mapper

import (
	"fmt"
	"strings"

	dg "github.com/bwmarrin/discordgo"
	"github.com/superloach/mapper/types"

	"github.com/go-snart/snart/bot"
	"github.com/go-snart/snart/db"
	"github.com/go-snart/snart/db/admin"
	"github.com/go-snart/snart/route"
)

// Searcher is a wrapper that creates a function suitable for route.
func Searcher(
	b *bot.Bot, filt func(*types.Location) bool,
	limit int, msg string,
) func(ctx *route.Ctx) error {
	return func(ctx *route.Ctx) error {
		return Search(
			b.DB, ctx, admin.IsAdmin(b.DB)(ctx),
			filt, limit, msg,
		)
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
	d *db.DB, ctx *route.Ctx, admin bool,
	filt func(*types.Location) bool, limit int, msg string,
) error {
	err := ctx.Session.ChannelTyping(ctx.Message.ChannelID)
	if err != nil {
		err = fmt.Errorf("typing %q: %w", ctx.Message.ChannelID, err)
		warn.Println(err)
	}

	debug := false

	if admin {
		ctx.Flag.BoolVar(&debug, "debug", false, "print extra info")
	}

	_ = ctx.Flag.Parse()

	args := ctx.Flag.Args()
	queries := cleanQueries(strings.Split(strings.Join(args, " "), "+"))

	if len(queries) == 0 {
		rep := ctx.Reply()
		rep.Content = fmt.Sprintf(
			"please specify a query.\nex: `%s%s [name of location]`",
			ctx.Prefix.Clean, ctx.Route.Name)

		return rep.Send()
	}

	locs, err := types.GetLocations(ctx, d)
	if err != nil {
		err = fmt.Errorf("get locs: %w", err)
		warn.Println(err)

		return err
	}

	for _, query := range queries {
		err = ctx.Session.ChannelTyping(ctx.Message.ChannelID)
		if err != nil {
			err = fmt.Errorf("typing %q: %w", ctx.Message.ChannelID, err)
			warn.Println(err)
		}

		searchQuery(
			ctx, query, locs,
			limit, debug, msg,
		)
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
	ctx *route.Ctx, query string, locations []*types.Location,
	limit int, _debug bool, msg string,
) {
	pss := search(clean(query), locations, 50, limit)

	pg := NewWidget(ctx.Session, ctx.Message.ChannelID, ctx.Message.Author.ID)

	pg.Add(firstEmbed(len(pss), msg, nick(ctx.Message)))

	for i, ps := range pss {
		debug.Println(ps)

		e := mkEmbed(ps, i, len(pss), msg, nick(ctx.Message))

		if _debug {
			debugEmbed(e, ps)
		}

		pg.Add(e)
	}

	go pg.Spawn()
}

func firstEmbed(total int, msg, name string) *dg.MessageEmbed {
	return &dg.MessageEmbed{
		Title: fmt.Sprintf("%d POIs Found", total),
		Description: fmt.Sprintf(
			"%s/%s: navigate left/right\n%s: confirm selection",
			EmoteLeft, EmoteRight, EmoteConfirm,
		),
		Footer: &dg.MessageEmbedFooter{
			Text: fmt.Sprintf(
				"%d %s • %s",
				total, msg, name),
		},
	}
}

func mkEmbed(ps *locationScore, i, pssl int, msg, name string) *dg.MessageEmbed {
	embed := &dg.MessageEmbed{
		Title: ps.Name,
		URL:   ps.URL(),
		Footer: &dg.MessageEmbedFooter{
			Text: fmt.Sprintf(
				"%d/%d %s • %s",
				i+1, pssl, msg, name),
		},
	}

	if ps.Image != nil {
		embed.Thumbnail = &dg.MessageEmbedThumbnail{
			URL:    *ps.Image,
			Height: 100,
		}
	}

	if ps.Notes != nil {
		addField(embed, "Notes", *ps.Notes)
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

func nick(m *dg.Message) string {
	if m.Member != nil {
		if m.Member.Nick != "" {
			return m.Member.Nick
		}

		if m.Member.User != nil {
			return m.Member.User.Username
		}
	}

	if m.Author != nil {
		return m.Author.Username
	}

	return "NAME UNKNOWN"
}
