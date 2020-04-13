package main

import (
	"fmt"
	"strings"

	dg "github.com/bwmarrin/discordgo"
	"github.com/go-snart/db"
	"github.com/go-snart/route"
	fuzzy "github.com/paul-mannino/go-fuzzywuzzy"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

func scorer(s1, s2 string) int {
	return (fuzzy.PartialRatio(s1, s2) + fuzzy.Ratio(s1, s2)) / 2
}

func clean(s string) string {
	return fuzzy.Cleanse(s, true)
}

func Poi(d *db.DB, ctx *route.Ctx) error {
	_f := "_poi"

	debug := ctx.Flags.Bool("debug", false, "print extra info")

	err := ctx.Flags.Parse()
	if err != nil {
		err = fmt.Errorf("flag parse: %w", err)
		Log.Error(_f, err)
		return err
	}

	args := ctx.Flags.Args()
	if len(args) == 0 {
		rep1 := ctx.Reply()
		rep1.Content = "please specify a query.\nex: `" +
			ctx.CleanPrefix + ctx.Route.Name + " name of poi`"

		_, err = rep1.Send()
		if err != nil {
			err = fmt.Errorf("rep1 send: %w", err)
			Log.Error(_f, err)
			return err
		}

		return nil
	}

	err = ctx.Session.ChannelTyping(ctx.Message.ChannelID)
	if err != nil {
		err = fmt.Errorf("typing %#v: %w", ctx.Message.ChannelID, err)
		Log.Error(_f, err)
		return err
	}

	query := strings.Join(args, " ")

	q := r.DB("poi").Table("poi")
	switch ctx.Route.Name {
	case "gyms":
		q = q.Filter(map[string]interface{}{
			"pkmn": "G",
		})
	case "stops":
		q = q.Filter(map[string]interface{}{
			"pkmn": "S",
		})
	default:
	}

	pois := make([]*POI, 0)
	err = q.ReadAll(&pois, d)
	if err != nil {
		err = fmt.Errorf("readall &pois: %w", err)
		Log.Error(_f, err)
		return err
	}

	choices := make([]string, 0)
	for _, poi := range pois {
		choices = append(choices, poi.Name)
		if poi.Alias != nil {
			for _, alias := range poi.Alias {
				choices = append(choices, alias)
			}
		}
	}

	suggNames, err := fuzzy.Extract(query, choices, 20, clean, scorer, 0)
	if err != nil {
		err = fmt.Errorf("fuzzy %#v: %w", query, err)
		Log.Error(_f, err)
		return err
	}
	if len(suggNames) == 0 {
		rep2 := ctx.Reply()
		rep2.Content = "no results found"

		_, err = rep2.Send()
		if err != nil {
			err = fmt.Errorf("rep2 send: %w", err)
			Log.Error(_f, err)
			return err
		}

		return nil
	}

	suggs := make([]*POI, 0)
	for _, sug := range suggNames {
		Log.Debugf(_f, "%#v\n", sug)
		for _, poi := range pois {
			if sug.Match == poi.Name {
				suggs = append(suggs, poi)
			}
			if poi.Alias != nil {
				for _, alias := range poi.Alias {
					if sug.Match == alias {
						suggs = append(suggs, poi)
					}
				}
			}
		}
	}

	pg := NewPager(ctx.Session, ctx.Message.ChannelID, ctx.Message.Author.ID)
	for i, sugg := range suggs {
		Log.Debugf(_f, "%#v\n", sugg)

		desc := []string{}
		if sugg.Notes != "" {
			desc = append(desc, "Notes: `"+sugg.Notes+"`")
		}
		if sugg.Alias != nil && len(sugg.Alias) > 0 {
			alias := strings.Join(sugg.Alias, "`, `")
			desc = append(desc, "Aliases: `"+alias+"`")
		}
		if *debug {
			desc = append(desc, "ID: `"+sugg.ID+"`")
			if sugg.Ingr != "" {
				desc = append(desc, "Ingress: `"+sugg.Ingr+"`")
			}
			if sugg.Pkmn != "" {
				desc = append(desc, "Pokemon: `"+sugg.Pkmn+"`")
			}
			if sugg.Wzrd != "" {
				desc = append(desc, "Wizards: `"+sugg.Wzrd+"`")
			}
		}

		embed := &dg.MessageEmbed{}
		embed.Title = sugg.Name
		embed.URL = sugg.URL()
		embed.Description = strings.Join(desc, "\n")
		embed.Thumbnail = &dg.MessageEmbedThumbnail{
			URL: sugg.Image,
		}
		embed.Footer = &dg.MessageEmbedFooter{
			Text: fmt.Sprintf("%d/%d", i+1, len(suggs)),
		}
		pg.Add(embed)
	}

	err = pg.Spawn()
	if err != nil {
		err = fmt.Errorf("pg spawn: %w", err)
		Log.Error(_f, err)
		return err
	}

	return nil
}
