package mapper

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	dg "github.com/bwmarrin/discordgo"
	"github.com/go-snart/snart/db"
	"github.com/go-snart/snart/route"
	fuzzy "github.com/paul-mannino/go-fuzzywuzzy"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

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
		if pss[i].score == pss[j].score {
			return pss[i].Name < pss[j].Name
		}
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

func Search(d *db.DB, ctx *route.Ctx, admin bool) error {
	_f := "Search"

	var debug *bool

	if admin {
		debug = ctx.Flags.Bool("debug", false, "print extra info")
	} else {
		_debug := false
		debug = &_debug
	}

	err := ctx.Flags.Parse()
	if err != nil {
		err = fmt.Errorf("flag parse: %w", err)
		Log.Error(_f, err)
		return err
	}

	args := ctx.Flags.Args()
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
			ctx.CleanPrefix + ctx.Route.Name + " name of poi`"

		return rep1.Send()
	}

	err = ctx.Session.ChannelTyping(ctx.Message.ChannelID)
	if err != nil {
		err = fmt.Errorf("typing %#v: %w", ctx.Message.ChannelID, err)
		Log.Warn(_f, err)
	}

	q := r.DB("mapper").Table("poi")

	bounds := GetBounds(d, ctx)
	Log.Debugf(_f, "bounds: %#v", bounds)

	if bounds != nil {
		q = q.Filter((*bounds).Intersects(r.Row.Field("loc")))
	}

	msg := "POIs"
	limit := 100

	switch ctx.Route.Name {
	case "gyms":
		q = q.Filter(map[string]interface{}{
			"pkmn": "G",
		})
		msg = "Gyms"
		limit = 25
	case "stops":
		q = q.Filter(map[string]interface{}{
			"pkmn": "S",
		})
		msg = "PokéStops"
		limit = 50
	default:
	}

	Log.Debugf(_f, "gonna readall %s", q)

	pois := make([]*POI, 0)
	err = q.ReadAll(&pois, d)
	if err != nil {
		err = fmt.Errorf("readall &pois: %w", err)
		Log.Error(_f, err)
		return err
	}

	Log.Debugf(_f, "read %d", len(pois))

	for _, query := range nqueries {
		err = ctx.Session.ChannelTyping(ctx.Message.ChannelID)
		if err != nil {
			err = fmt.Errorf("typing %#v: %w", ctx.Message.ChannelID, err)
			Log.Warn(_f, err)
		}

		err = searchQuery(
			d, ctx,
			query, pois, limit,
			*debug, msg,
		)
		if err != nil {
			Log.Warn(_f, err)
		}
	}

	return nil
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

		embed := &dg.MessageEmbed{
			Title: ps.Name,
			Description: fmt.Sprintf(
				"%s suggested for %s",
				msg, nick(ctx.Message),
			),
			URL: ps.URL(),
			Thumbnail: &dg.MessageEmbedThumbnail{
				URL: ps.Image,
			},
			Footer: &dg.MessageEmbedFooter{
				Text: fmt.Sprintf("%d/%d", i+1, len(pss)),
			},
		}

		if ps.Notes != "" {
			embed.Fields = append(embed.Fields,
				&dg.MessageEmbedField{
					Name:   "Notes",
					Value:  ps.Notes,
					Inline: false,
				},
			)
		}

		if ps.Alias != nil && len(ps.Alias) > 0 {
			embed.Fields = append(embed.Fields,
				&dg.MessageEmbedField{
					Name:   "Aliases",
					Value:  strings.Join(ps.Alias, "\n"),
					Inline: false,
				},
			)
		}

		if debug {
			dmsg := []string{}

			dmsg = append(dmsg, "ID: `"+ps.ID+"`")
			if ps.Ingr != "" {
				dmsg = append(dmsg, "Ingress: `"+ps.Ingr+"`")
			}
			if ps.Pkmn != "" {
				dmsg = append(dmsg, "Pokemon: `"+ps.Pkmn+"`")
			}
			if ps.Wzrd != "" {
				dmsg = append(dmsg, "Wizards: `"+ps.Wzrd+"`")
			}
			dmsg = append(dmsg, "Score: "+strconv.Itoa(ps.score)+"%")
			dmsg = append(dmsg, "Link: "+ps.URL())

			embed.Fields = append(embed.Fields,
				&dg.MessageEmbedField{
					Name:   "Aliases",
					Value:  strings.Join(dmsg, "\n"),
					Inline: false,
				},
			)
		}

		pg.Add(embed)
	}

	go pg.Spawn()

	return nil
}
