package mapper

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	dg "github.com/bwmarrin/discordgo"
	"github.com/go-snart/db"
	"github.com/go-snart/route"
	fuzzy "github.com/paul-mannino/go-fuzzywuzzy"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

func words(s1, s2 string) int {
	a := 0
	s1s := strings.Split(s1, " ")
	s2s := strings.Split(s2, " ")
	for _, w1 := range s1s {
		for _, w2 := range s2s {
			if w1 == w2 {
				a++
				break
			}
		}
	}
	a *= 100
	a /= len(s1s)
	return a
}

func scorer(s1, s2 string) int {
	score := 0
	score += 3 * fuzzy.PartialRatio(s1, s2)
	score += 2 * fuzzy.Ratio(s1, s2)
	score += 1 * words(s1, s2)
	return score / 6
}

type poiScore struct {
	*POI
	score int
}

func score(q string, p *POI) *poiScore {
	ps := &poiScore{}

	ps.POI = p
	ps.score = scorer(q, clean(p.Name))

	for _, a := range p.Alias {
		s := scorer(q, clean(a))
		if s > ps.score {
			ps.score = s
		}
	}

	return ps
}

func search(q string, ps []*POI, min, num int) []*poiScore {
	pss := make([]*poiScore, len(ps))
	for i, p := range ps {
		pss[i] = score(q, p)
	}

	sort.Slice(pss, func(i, j int) bool {
		return pss[i].score > pss[j].score
	})

	i := 0
	for ; i < num && i < len(pss); i++ {
		if pss[i].score < min {
			break
		}
	}

	return pss[:i]
}

func clean(s string) string {
	return fuzzy.Cleanse(s, true)
}

func Search(d *db.DB, ctx *route.Ctx) error {
	_f := "Search"

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

		return rep1.Send()
	}

	for _, query := range strings.Split(strings.Join(args, " "), "+") {
		err = ctx.Session.ChannelTyping(ctx.Message.ChannelID)
		if err != nil {
			err = fmt.Errorf("typing %#v: %w", ctx.Message.ChannelID, err)
			Log.Error(_f, err)
			return err
		}

		err = searchQuery(d, ctx, query, *debug)
		if err != nil {
			Log.Warn(_f, err)
		}
	}

	return nil
}

func searchQuery(d *db.DB, ctx *route.Ctx, query string, debug bool) error {
	_f := "searchQuery"

	q := r.DB("mapper").Table("poi")
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
	err := q.ReadAll(&pois, d)
	if err != nil {
		err = fmt.Errorf("readall &pois: %w", err)
		Log.Error(_f, err)
		return err
	}

	pss := search(clean(query), pois, 50, 20)
	if len(pss) == 0 {
		rep2 := ctx.Reply()
		rep2.Content = "no results found"

		return rep2.Send()
	}

	pg := NewWidget(ctx.Session, ctx.Message.ChannelID, ctx.Message.Author.ID)
	for i, ps := range pss {
		Log.Debugf(_f, "%#v\n", ps)

		desc := []string{}
		if ps.Notes != "" {
			desc = append(desc, "Notes:\n```\n"+ps.Notes+"\n```")
		}
		if ps.Alias != nil && len(ps.Alias) > 0 {
			alias := "`" + strings.Join(ps.Alias, "`, `") + "`"
			desc = append(desc, "Aliases: "+alias)
		}
		if debug {
			desc = append(desc, "ID: `"+ps.ID+"`")
			if ps.Ingr != "" {
				desc = append(desc, "Ingress: `"+ps.Ingr+"`")
			}
			if ps.Pkmn != "" {
				desc = append(desc, "Pokemon: `"+ps.Pkmn+"`")
			}
			if ps.Wzrd != "" {
				desc = append(desc, "Wizards: `"+ps.Wzrd+"`")
			}
			desc = append(desc, "Score: "+strconv.Itoa(ps.score)+"%")
		}

		embed := &dg.MessageEmbed{}
		embed.Title = ps.Name
		embed.URL = ps.URL()
		embed.Description = strings.Join(desc, "\n")
		embed.Thumbnail = &dg.MessageEmbedThumbnail{
			URL: ps.Image,
		}
		embed.Footer = &dg.MessageEmbedFooter{
			Text: fmt.Sprintf("%d/%d", i+1, len(pss)),
		}
		pg.Add(embed)
	}

	go pg.Spawn()

	return nil
}
