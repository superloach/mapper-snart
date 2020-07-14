package mapper

import (
	"fmt"

	"github.com/go-snart/snart/db"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
	"gopkg.in/rethinkdb/rethinkdb-go.v6/types"
)

// NeighTable is a table builder for mapper.neigh.
var NeighTable = db.BuildTable(MapperDB, "neigh")

// Neigh contains the ID, Name, and center point of a neighborhood.
type Neigh struct {
	ID   string       `rethinkdb:"id"`
	Name string       `rethinkdb:"name"`
	Loc  *types.Point `rethinkdb:"loc"`
}

// NeighCache maintains a running state of known Neighs.
func NeighCache(d *db.DB) {
	_f := "NeighCache"

	d.WaitReady()

	q := NeighTable.Build(d).Changes(
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
		New *Neigh `rethinkdb:"new_val"`
		Old *Neigh `rethinkdb:"old_val"`
	}{}

	d.Cache.Lock()
	if !d.Cache.Has("mapper.neigh") {
		d.Cache.Set("mapper.neigh", db.NewMapCache())
	}

	cache := d.Cache.Get("mapper.neigh").(db.Cache)
	d.Cache.Unlock()

	for curs.Next(&chng) {
		cache.Lock()

		if chng.New != nil {
			cache.Set(chng.New.ID, chng.New)
		} else {
			cache.Del(chng.Old.ID)
		}

		cache.Unlock()

		clearLocationNeighs(d)
	}

	if err := curs.Err(); err != nil {
		resp, ok := curs.NextResponse()

		err = fmt.Errorf(
			"curs err: %w\n"+
				"chng is %#v/%#v\n"+
				"resp(%v) is %q",
			err,
			chng.New, chng.Old,
			ok, resp,
		)
		Log.Error(_f, err)

		return
	}
}

func clearLocationNeighs(d *db.DB) {
	d.Cache.Lock()
	locationCache := d.Cache.Get("mapper.location").(db.Cache)
	d.Cache.Unlock()

	locationCache.Lock()
	keys := locationCache.Keys()
	locationCache.Unlock()

	for _, k := range keys {
		locationCache.Lock()
		location := locationCache.Get(k).(*Location)
		locationCache.Unlock()

		location.Neigh = nil
	}
}
