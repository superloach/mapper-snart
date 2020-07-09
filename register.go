package mapper

import (
	"github.com/superloach/minori"

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

	go POICache(b.DB)
	go BoundCache(b.DB)
	go NeighCache(b.DB)

	b.AddGamer(bot.GamerText(
		"Pokémon GO",
		dg.GameTypeGame,
	))
	b.AddGamer(GamerCounts(
		"%.f Gyms | %.f PokéStops",
		map[string]interface{}{"pkmn": "G"},
		map[string]interface{}{"pkmn": "S"},
	))

	b.AddGamer(bot.GamerText(
		"Ingress",
		dg.GameTypeGame,
	))
	b.AddGamer(GamerCounts(
		"%.f POIs",
		map[string]interface{}{},
	))

	registerCmds(b)
	registerAdminCmds(b)

	Log.Info(_f, "registered")

	return nil
}

func registerCmds(b *bot.Bot) {
	b.Router.Add(
		&route.Route{
			Name:  "poi",
			Match: "pois?",
			Desc:  "Search for any POIs.",
			Cat:   name,
			Okay:  nil,
			Func:  Searcher(b, 0),
		},
		&route.Route{
			Name:  "g",
			Match: "(g(yms?)?)|(fighty place)",
			Desc:  "Search for Pokemon Go gyms.",
			Cat:   name,
			Okay:  nil,
			Func:  Searcher(b, 'g'),
		},
		&route.Route{
			Name:  "p",
			Match: "(p(ok[eé]stops?)?)|(s(tops?)?)",
			Desc:  "Search for Pokemon Go stops.",
			Cat:   name,
			Okay:  nil,
			Func:  Searcher(b, 'p'),
		},
	)
}

func registerAdminCmds(b *bot.Bot) {
	b.Router.Add(
		&route.Route{
			Name:  "count",
			Match: "counts?",
			Desc:  "get poi counts",
			Cat:   name,
			Okay:  b.Admin,
			Func:  Counts,
		},
	)
}
