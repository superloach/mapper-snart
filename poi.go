package mapper

import (
	"fmt"

	"github.com/go-snart/snart/db"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
	"gopkg.in/rethinkdb/rethinkdb-go.v6/types"
)

// POITable is a table builder for mapper.poi.
var POITable = db.BuildTable(
	MapperDB, "poi",
	nil, nil,
)

// POI contains lots of info about a Point Of Interest in location-based games.
type POI struct {
	ID    string `rethinkdb:"id"`
	Name  string `rethinkdb:"name"`
	Image string `rethinkdb:"image"`
	Notes string `rethinkdb:"notes"`

	Loc *types.Point `rethinkdb:"loc"`

	Ingr string `rethinkdb:"ingr"`
	Pkmn string `rethinkdb:"pkmn"`
	Wzrd string `rethinkdb:"wzrd"`

	Alias []string `rethinkdb:"alias"`

	Neigh *string `rethinkdb:"-"`
}

// URL returns a suitable URL for finding directions to the POI.
func (p *POI) URL() string {
	return mapURL(fmt.Sprintf(
		"%.06f,%.06f",
		p.Loc.Lat,
		p.Loc.Lon,
	))
}

// POICache maintains a running state of known POIs.
func POICache(d *db.DB) {
	_f := "POICache"

	d.WaitReady()

	if !d.Cache.Has("mapper.poi") {
		d.Cache.Lock()
		d.Cache.Set("mapper.poi", db.NewMapCache())
		d.Cache.Unlock()
	}

	q := POITable.Changes(
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
		New *POI `rethinkdb:"new_val"`
		Old *POI `rethinkdb:"old_val"`
	}{}

	d.Cache.Lock()
	cache, _ := d.Cache.Get("mapper.poi").(db.Cache)
	d.Cache.Unlock()

	for curs.Next(&chng) {
		cache.Lock()

		switch {
		case chng.New != nil:
			cache.Set(chng.New.ID, chng.New)
		case chng.Old != nil:
			cache.Del(chng.Old.ID)
		}

		cache.Unlock()
	}

	if err := curs.Err(); err != nil {
		err = fmt.Errorf("curs err: %w", err)
		Log.Error(_f, err)

		return
	}
}
