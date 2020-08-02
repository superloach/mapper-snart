// Package mapper is a Snart plugin which provides mapping for Locations in various location-based games.
package mapper

import (
	"github.com/go-snart/snart/bot"
	"github.com/go-snart/snart/logs"
	"github.com/go-snart/snart/route"

	"github.com/superloach/mapper/types"
)

const _p = "mapper"

var debug, _, warn = logs.Loggers(_p)

func init() {
	bot.Register(_p, Register)
}

// Register adds all of mapper's handlers to the given bot.
func Register(b *bot.Bot) error {
	const _f = "Register"

	debug.Println(_f, "registering")

	registerCmds(b)

	debug.Println(_f, "registered")

	return nil
}

func registerCmds(b *bot.Bot) {
	b.Router.Add(
		&route.Route{
			Name:  "poi",
			Match: "pois?",
			Desc:  "Search for Niantic POIs.",
			Cat:   _p,
			Okay:  nil,
			Func: Searcher(b, func(l *types.Location) bool {
				return true
			}, 100, "POIs"),
		},

		&route.Route{
			Name:  "p",
			Match: "(p(ok[eé]stops?)?)|(s(tops?)?)",
			Desc:  "Search for Pokémon GO Stops.",
			Cat:   _p,
			Okay:  nil,
			Func: Searcher(b, func(l *types.Location) bool {
				return l.PkmnType == types.PkmnTypeStop
			}, 50, "PokéStops"),
		},

		&route.Route{
			Name:  "g",
			Match: "(g(yms?)?)|(fighty place)",
			Desc:  "Search for Pokémon GO Gyms.",
			Cat:   _p,
			Okay:  nil,
			Func: Searcher(b, func(l *types.Location) bool {
				return l.PkmnType == types.PkmnTypeGym ||
					l.PkmnType == types.PkmnTypeEXGym
			}, 25, "Gyms"),
		},
	)
}
