package mapper

import (
	"github.com/go-snart/snart/route"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

func Qtest(c *route.Ctx) (r.Term, error) {
	return r.DB("mapper").Table("poi").Count(), nil
}
