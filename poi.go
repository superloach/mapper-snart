package main

import (
	"strings"

	dgw "github.com/Necroforger/dgwidgets"
	dg "github.com/bwmarrin/discordgo"
	"github.com/paul-mannino/go-fuzzywuzzy"
	"github.com/go-snart/snart/lib/db"
	"github.com/go-snart/snart/lib/errs"
	"github.com/go-snart/snart/lib/route"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

func scorer(s1, s2 string) int {
	return (fuzzy.PartialRatio(s1, s2) + fuzzy.Ratio(s1, s2)) / 2
}

func clean(s string) string {
	return fuzzy.Cleanse(s, false)
}

func Poi(d *db.DB) func(ctx *route.Ctx) error {
	_f := "Poi"
	return func(ctx *route.Ctx) error {
		err := _poi(d, ctx)
		if err != nil {
			errs.Wrap(&err, `_poi(d, ctx)`)
			Log.Error(_f, err)
			return err
		}

		return nil
	}
}

func _poi(d *db.DB, ctx *route.Ctx) error {
	_f := "_poi"

	err := ctx.Flags.Parse()
	if err != nil {
		errs.Wrap(&err, `ctx.Flags.Parse()`)
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
			errs.Wrap(&err, `rep1.Send()`)
			Log.Error(_f, err)
			return err
		}

		return nil
	}

	query := strings.Join(args, " ")

	pois := make([]*POI, 0)
	q := r.DB("poi").Table("poi")
	err = q.ReadAll(&pois, d)
	if err != nil {
		errs.Wrap(&err, `q.ReadAll(&pois, d)`)
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

	suggNames, err := fuzzy.Extract(query, choices, 12, clean, scorer, 50)
	if err != nil {
		errs.Wrap(&err, `fuzzy.Extract(%#v, choices, 12, clean, scorer, 50)`, query)
		Log.Error(_f, err)
		return err
	}
	if len(suggNames) == 0 {
		rep2 := ctx.Reply()
		rep2.Content = "no results found"

		_, err = rep2.Send()
		if err != nil {
			errs.Wrap(&err, `rep2.Send()`)
			Log.Error(_f, err)
			return err
		}

		return nil
	}

	suggs := make([]*POI, 0)
	for _, poi := range pois {
		for _, sug := range suggNames {
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

	pg := dgw.NewPaginator(ctx.Session, ctx.Message.ChannelID)
	for _, sugg := range suggs {
		embed := &dg.MessageEmbed{}
		embed.Title = sugg.Name
		embed.URL = MapURL(sugg)
		embed.Description = sugg.Notes
		embed.Thumbnail = &dg.MessageEmbedThumbnail{
			URL: sugg.Image,
		}
		pg.Add(embed)
	}

	err = pg.Spawn()
	if err != nil {
		errs.Wrap(&err, `pg.Spawn()`)
		Log.Error(_f, err)
		return err
	}

	return nil
}
