package mapper

import (
	"fmt"

	"github.com/go-snart/snart/db"
	"github.com/go-snart/snart/route"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
	"gopkg.in/rethinkdb/rethinkdb-go.v6/types"
)

type Bounds struct {
	ID    string      `json:"id" rethinkdb:"id"`
	Value types.Lines `json:"value" rethinkdb:"value"`
}

func GetBounds(d *db.DB, ctx *route.Ctx) *r.Term {
	_f := "GetBounds"

	ids := r.Expr([]interface{}{
		ctx.Message.ChannelID,
		ctx.Message.GuildID,
	})

	q := r.DB("mapper").Table("bounds").Filter(
		func(row r.Term) r.Term {
			return ids.Contains(row.Field("id"))
		},
	)

	boundss := make([]*Bounds, 0)
	err := q.ReadAll(&boundss, d)
	if err != nil {
		err = fmt.Errorf("readall &boundss: %w", err)
		Log.Warn(_f, err)
		return nil
	}

	val := types.Lines(nil)

	for _, bounds := range boundss {
		if bounds.ID == ctx.Message.ChannelID {
			val = bounds.Value
		}
	}

	if val == nil {
		for _, bounds := range boundss {
			if bounds.ID == ctx.Message.GuildID {
				val = bounds.Value
			}
		}
	}

	if val != nil {
		rql, err := val.MarshalRQL()
		if err != nil {
			err = fmt.Errorf("val rql: %w", err)
			Log.Warn(_f, err)
			return nil
		}

		bounds := r.Expr(rql)
		return &bounds
	}

	return nil
}
