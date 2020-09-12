package main

import (
	"context"

	"github.com/go-snart/snart/bot/plug"
	"github.com/go-snart/snart/db"
	"github.com/go-snart/snart/route"

	"github.com/superloach/mapper/types"
)

var Plug = plug.Plug(&Mapper{})

func main() {}

type Mapper struct {
	*plug.Base
	DB *db.DB
}

func (m *Mapper) String() string {
	return "mapper"
}

func (m *Mapper) PlugDB(d *db.DB) {
	m.Base.PlugDB(d)

	m.DB = db.New(context.Background(), "mapper")
}

func (m *Mapper) PlugHandler(r *route.Handler) {
	m.Base.PlugHandler(r)

	r.Add(
		&route.Route{
			Name:  "poi",
			Match: route.MustMatch("pois?"),
			Desc:  "Search for Niantic POIs.",
			Cat:   m.String(),
			Okay:  nil,
			Func: m.Searcher(func(l *types.Location) bool {
				return true
			}, 100, "POIs"),
		},

		&route.Route{
			Name:  "p",
			Match: route.MustMatch("(p(ok[eé]stops?)?)|(s(tops?)?)"),
			Desc:  "Search for Pokémon GO Stops.",
			Cat:   m.String(),
			Okay:  nil,
			Func: m.Searcher(func(l *types.Location) bool {
				return l.PkmnType == types.PkmnTypeStop
			}, 50, "PokéStops"),
		},

		&route.Route{
			Name:  "g",
			Match: route.MustMatch("(g(yms?)?)|(fighty place)"),
			Desc:  "Search for Pokémon GO Gyms.",
			Cat:   m.String(),
			Okay:  nil,
			Func: m.Searcher(func(l *types.Location) bool {
				return l.PkmnType == types.PkmnTypeGym ||
					l.PkmnType == types.PkmnTypeEXGym
			}, 25, "Gyms"),
		},
	)
}
