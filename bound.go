package mapper

import (
	"fmt"

	"github.com/go-snart/snart/db"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
	"gopkg.in/rethinkdb/rethinkdb-go.v6/types"
)

// BoundTable is a table builder for mapper.bound.
var BoundTable = db.BuildTable(MapperDB, "bound")

// Bound contains the ID and border of a Guild/Channel boundary.
type Bound struct {
	ID    string      `rethinkdb:"id"`
	Value types.Lines `rethinkdb:"value"`
}

// BoundCache maintains a running state of known Bounds.
func BoundCache(d *db.DB) {
	_f := "BoundCache"

	d.WaitReady()

	q := BoundTable.Build(d).Changes(
		r.ChangesOpts{IncludeInitial: true},
	)

	curs, err := q.Run(d)
	if err != nil {
		err = fmt.Errorf("db run %s: %w", q, err)
		Log.Error(_f, err)

		return
	}
	defer curs.Close()

	chng := struct {
		New *Bound `rethinkdb:"new_val"`
		Old *Bound `rethinkdb:"old_val"`
	}{}

	for curs.Next(&chng) {
		d.Cache.Lock()

		if chng.New != nil {
			d.Cache.Set("mapper.bound."+chng.New.ID, db.NewMapCache())

			go BoundLocationCache(d, chng.New)
		} else {
			d.Cache.Del("mapper.bound." + chng.Old.ID)
		}

		d.Cache.Unlock()
	}

	if err := curs.Err(); err != nil {
		err = fmt.Errorf("curs err: %w", err)
		Log.Error(_f, err)

		return
	}
}

// BoundLocationCache maintains a running state of Locations within a Bound.
func BoundLocationCache(d *db.DB, bound *Bound) {
	_f := "BoundLocationCache"

	d.WaitReady()

	rql, err := bound.Value.MarshalRQL()
	if err != nil {
		err = fmt.Errorf("val rql: %w", err)
		Log.Error(_f, err)

		return
	}

	val := r.Expr(rql)

	q := LocationTable.Build(d).Filter(
		val.Intersects(r.Row.Field("loc")),
	).Field("id").Changes(
		r.ChangesOpts{IncludeInitial: true},
	)

	curs, err := q.Run(d)
	if err != nil {
		err = fmt.Errorf("db run %s: %w", q, err)
		Log.Error(_f, err)

		return
	}
	defer curs.Close()

	chng := struct {
		New *string `rethinkdb:"new_val"`
		Old *string `rethinkdb:"old_val"`
	}{}

	d.Cache.Lock()
	cache := d.Cache.Get("mapper.bound." + bound.ID).(db.Cache)
	d.Cache.Unlock()

	for curs.Next(&chng) {
		cache.Lock()

		if chng.New != nil {
			cache.Set(*chng.New, true)
		} else {
			cache.Del(*chng.Old)
		}

		cache.Unlock()
	}

	if err := curs.Err(); err != nil {
		err = fmt.Errorf("curs err: %w", err)
		Log.Error(_f, err)

		return
	}
}
