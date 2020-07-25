package mapper

import (
	"context"
	"fmt"

	"github.com/dewski/spatial"

	"github.com/go-snart/snart/db"
)

// LocationTable is a table builder for mapper.location.
func LocationTable(ctx context.Context, d *db.DB) {
	const q = `CREATE TABLE IF NOT EXISTS location(
		id TEXT PRIMARY KEY UNIQUE NOT NULL,
		name TEXT,
		image TEXT,
		notes TEXT,

		value GEOMETRY(POINT),

		ingrtype INT,
		pkmntype INT,
		wzrdtype INT,

		aliases TEXT[]
	)`
	x, err := d.Conn(&ctx).Exec(ctx, q)
	Log.Debugf("locationtable", "%#v %#v", x, err)
}

// Location contains lots of info about a Point Of Interest in location-based games.
type Location struct {
	ID    string
	Name  string
	Image string
	Notes string

	Value spatial.Point

	IngrType IngrType
	PkmnType PkmnType
	WzrdType WzrdType

	Aliases []string
}

// URL returns a suitable URL for finding directions to the Location.
func (p *Location) URL() string {
	return mapURL(fmt.Sprintf(
		"%.06f,%.06f",
		p.Value.Lat,
		p.Value.Lng,
	))
}

func GetLocations(d *db.DB, ctx context.Context) ([]*Location, error) {
	_f := "GetLocations"

	const q = `
		SELECT
			id, name, image, notes,
			value,
			ingrtype, pkmntype, wzrdtype,
			aliases
		FROM location
	`

	rows, err := d.Conn(&ctx).Query(ctx, q)
	if err != nil {
		err = fmt.Errorf("db run %s: %w", q, err)
		Log.Error(_f, err)

		return nil, err
	}
	defer rows.Close()

	locs := []*Location(nil)

	for rows.Next() {
		loc := &Location{}

		err := rows.Scan(
			&loc.ID, &loc.Name, &loc.Image, &loc.Notes,
			&loc.Value,
			&loc.IngrType, &loc.PkmnType, &loc.WzrdType,
			&loc.Aliases,
		)
		if err != nil {
			return nil, err
		}

		locs = append(locs, loc)
	}

	err = rows.Err()
	if err != nil {
		err = fmt.Errorf("curs err: %w", err)
		Log.Error(_f, err)

		return nil, err
	}

	return locs, nil
}
