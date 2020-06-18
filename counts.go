package mapper

import (
	"github.com/go-snart/snart/route"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

func Counts(c *route.Ctx) (r.Term, error) {
	return r.Expr(map[string]interface{}{
		"pois": r.DB("mapper").Table("poi").Count(),
		"gyms": r.DB("mapper").Table("poi").Filter(map[string]string{
			"pkmn": "G",
		}).Count(),
		"pokestops": r.DB("mapper").Table("poi").Filter(map[string]string{
			"pkmn": "S",
		}).Count(),
	}), nil
}
