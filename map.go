package main

import (
	"fmt"
	"net/url"
	"strings"

	dg "github.com/bwmarrin/discordgo"
	"github.com/go-snart/snart/log"
	"github.com/go-snart/snart/route"
)

// Map allows a user to get directions to an arbitrary location.
func Map(ctx *route.Ctx) error {
	const _f = "Map"

	_ = ctx.Flag.Parse()

	args := ctx.Flag.Args()
	query := strings.Join(args, " ")
	queries := strings.Split(query, "+")
	nqueries := make([]string, 0)

	for _, q := range queries {
		q = strings.TrimSpace(q)
		if len(q) == 0 {
			continue
		}

		nqueries = append(nqueries, q)
	}

	if len(nqueries) == 0 {
		rep1 := ctx.Reply()
		rep1.Content = "please specify a query.\nex: `" +
			ctx.Prefix.Clean + ctx.Route.Name + " name of place`"

		return rep1.Send()
	}

	msg := "Map given for"

	for _, query := range nqueries {
		err := ctx.Session.ChannelTyping(ctx.Message.ChannelID)
		if err != nil {
			err = fmt.Errorf("typing %q: %w", ctx.Message.ChannelID, err)
			log.Warn.Println(_f, err)
		}

		rep := ctx.Reply()
		rep.Embed = &dg.MessageEmbed{
			Title: query,
			URL:   MapURL(query),
			Footer: &dg.MessageEmbedFooter{
				Text: fmt.Sprintf(
					"%s %s",
					msg, nick(ctx.Message),
				),
			},
		}

		err = rep.Send()
		if err != nil {
			log.Warn.Println(_f, err)
		}
	}

	return nil
}

func MapURL(s string) string {
	s = url.PathEscape(s)
	s = strings.ReplaceAll(s, url.PathEscape(","), ",")

	return "https://www.google.com/maps/dir//" + s
}
