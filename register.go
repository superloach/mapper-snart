// Package mapper is a Snart plugin which provides mapping for Locations in various location-based games.
package mapper

import (
	"github.com/superloach/minori"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"

	dg "github.com/bwmarrin/discordgo"
	"github.com/go-snart/snart/bot"
	"github.com/go-snart/snart/db"
	"github.com/go-snart/snart/route"
)

const name = "mapper"

// Log is the logger for mapper.
var Log = minori.GetLogger(name)

// MapperDB is a DB builder for mapper.
var MapperDB = db.BuildDB("mapper")

func init() {
	bot.Register(name, Register)
}

// Register adds all of mapper's handlers to the given bot.
func Register(b *bot.Bot) error {
	_f := "Register"
	Log.Info(_f, "registering")

	go LocationCache(b.DB)
	go BoundCache(b.DB)
	go NeighCache(b.DB)

	registerGamers(b)
	registerCmds(b)
	registerAdminCmds(b)

	Log.Info(_f, "registered")

	return nil
}

func registerGamers(b *bot.Bot) {
	b.AddGamer(bot.GamerText("Pokémon GO", dg.GameTypeGame))
	b.AddGamer(GamerCounts("%.f Gyms | %.f PokéStops",
		r.Row.Field("pkmn").Eq(PkmnTypeGym).Or(
			r.Row.Field("pkmn").Eq(PkmnTypeEXGym)),
		r.Row.Field("pkmn").Eq(PkmnTypeStop)))

	b.AddGamer(bot.GamerText("Ingress", dg.GameTypeGame))
	b.AddGamer(GamerCounts("%.f Portals",
		r.Row.Field("ingr").Eq(IngrTypePortal)))
}

func registerCmds(b *bot.Bot) {
	b.Router.Add(&route.Route{
		Name:  "poi",
		Match: "pois?",
		Desc:  "Search for Niantic POIs.",
		Cat:   name,
		Okay:  nil,
		Func: Searcher(b, func(p *Location) bool {
			return true
		}, 100, "POIs"),
	}, &route.Route{
		Name:  "p",
		Match: "(p(ok[eé]stops?)?)|(s(tops?)?)",
		Desc:  "Search for Pokémon GO Stops.",
		Cat:   name,
		Okay:  nil,
		Func: Searcher(b, func(p *Location) bool {
			return p.PkmnType == PkmnTypeStop
		}, 50, "PokéStops"),
	}, &route.Route{
		Name:  "g",
		Match: "(g(yms?)?)|(fighty place)",
		Desc:  "Search for Pokémon GO Gyms.",
		Cat:   name,
		Okay:  nil,
		Func: Searcher(b, func(p *Location) bool {
			return p.PkmnType == PkmnTypeGym ||
				p.PkmnType == PkmnTypeEXGym
		}, 25, "Gyms"),
	})
}

func registerAdminCmds(b *bot.Bot) {
	b.Router.Add(&route.Route{
		Name:  "count",
		Match: "counts?",
		Desc:  "get poi counts",
		Cat:   name,
		Okay:  b.DB.Admin,
		Func:  Counts,
	})
}
