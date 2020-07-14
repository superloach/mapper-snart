package mapper

import (
	"fmt"

	"github.com/go-snart/snart/db"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
	"gopkg.in/rethinkdb/rethinkdb-go.v6/types"
)

// LocationTable is a table builder for mapper.location.
var LocationTable = db.BuildTable(MapperDB, "location")

// Location contains lots of info about a Point Of Interest in location-based games.
type Location struct {
	ID    string            `rethinkdb:"id"`
	Name  string            `rethinkdb:"name"`
	Image string            `rethinkdb:"image"`
	Notes map[string]string `rethinkdb:"notes"`

	Loc *types.Point `rethinkdb:"loc"`

	IngrType IngrType `rethinkdb:"ingr"`
	PkmnType PkmnType `rethinkdb:"pkmn"`
	WzrdType WzrdType `rethinkdb:"wzrd"`

	Alias []string `rethinkdb:"alias"`

	Neigh *string `rethinkdb:"-"`
}

// URL returns a suitable URL for finding directions to the Location.
func (p *Location) URL() string {
	return mapURL(fmt.Sprintf(
		"%.06f,%.06f",
		p.Loc.Lat,
		p.Loc.Lon,
	))
}

// LocationCache maintains a running state of known Locations.
func LocationCache(d *db.DB) {
	_f := "LocationCache"

	d.WaitReady()

	if !d.Cache.Has("mapper.location") {
		d.Cache.Lock()
		d.Cache.Set("mapper.location", db.NewMapCache())
		d.Cache.Unlock()
	}

	q := LocationTable.Build(d).Changes(
		r.ChangesOpts{IncludeInitial: true})

	curs, err := q.Run(d)
	if err != nil {
		err = fmt.Errorf("db run %s: %w", q, err)
		Log.Error(_f, err)

		return
	}
	defer curs.Close()

	chng := struct {
		New *Location `rethinkdb:"new_val"`
		Old *Location `rethinkdb:"old_val"`
	}{}

	d.Cache.Lock()
	cache, _ := d.Cache.Get("mapper.location").(db.Cache)
	d.Cache.Unlock()

	for curs.Next(&chng) {
		cache.Lock()

		if chng.New == nil {
			cache.Del(chng.Old.ID)
		} else {
			cache.Set(chng.New.ID, chng.New)
		}

		cache.Unlock()
	}

	if err := curs.Err(); err != nil {
		err = fmt.Errorf("curs err: %w", err)
		Log.Error(_f, err)

		return
	}
}
