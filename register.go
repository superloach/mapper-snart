// Package mapper is a Snart plugin which provides mapping for Locations in various location-based games.
package mapper

import (
	"github.com/superloach/minori"

	"github.com/go-snart/snart/bot"
	"github.com/go-snart/snart/route"
)

const name = "mapper"

// Log is the logger for mapper.
var Log = minori.GetLogger(name)

func init() {
	bot.Register(name, Register)
}

// Register adds all of mapper's handlers to the given bot.
func Register(b *bot.Bot) error {
	_f := "Register"
	Log.Info(_f, "registering")

	registerCmds(b)

	Log.Info(_f, "registered")

	return nil
}

func registerCmds(b *bot.Bot) {
	b.Router.Add(
		&route.Route{
			Name:  "poi",
			Match: "pois?",
			Desc:  "Search for Niantic POIs.",
			Cat:   name,
			Okay:  nil,
			Func: Searcher(b, func(p *Location) bool {
				return true
			}, 100, "POIs"),
		},

		&route.Route{
			Name:  "p",
			Match: "(p(ok[eé]stops?)?)|(s(tops?)?)",
			Desc:  "Search for Pokémon GO Stops.",
			Cat:   name,
			Okay:  nil,
			Func: Searcher(b, func(p *Location) bool {
				return p.PkmnType == PkmnTypeStop
			}, 50, "PokéStops"),
		},

		&route.Route{
			Name:  "g",
			Match: "(g(yms?)?)|(fighty place)",
			Desc:  "Search for Pokémon GO Gyms.",
			Cat:   name,
			Okay:  nil,
			Func: Searcher(b, func(p *Location) bool {
				return p.PkmnType == PkmnTypeGym ||
					p.PkmnType == PkmnTypeEXGym
			}, 25, "Gyms"),
		},
	)
}
