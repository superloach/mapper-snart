package types

import (
	"context"
	"fmt"

	"github.com/dewski/spatial"

	"github.com/go-snart/snart/db"
)

// LocationTable is a table builder for mapper.location.
func LocationTable(ctx context.Context, d *db.DB) {
	const (
		_f = "LocationTable"
		e  = `CREATE TABLE IF NOT EXISTS location(
			id TEXT PRIMARY KEY UNIQUE NOT NULL,
			name TEXT NOT NULL,
			image TEXT,
			notes TEXT,

			value GEOMETRY(POINT) NOT NULL,

			ingrtype INT NOT NULL,
			pkmntype INT NOT NULL,
			wzrdtype INT NOT NULL,

			aliases TEXT[]
		)`
	)

	_, err := d.Conn(&ctx).Exec(ctx, e)
	if err != nil {
		err = fmt.Errorf("exec %#q: %w", e, err)

		warn.Println(_f, err)

		return
	}
}

// Location contains lots of info about a Point Of Interest in location-based games.
type Location struct {
	ID    string
	Name  string
	Image *string
	Notes *string

	Value spatial.Point

	IngrType IngrType
	PkmnType PkmnType
	WzrdType WzrdType

	Aliases []string
}

// GetLocations retrieves a list of Locations from the given DB.
func GetLocations(ctx context.Context, d *db.DB) ([]*Location, error) {
	const _f = "GetLocations"

	LocationTable(ctx, d)

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
		err = fmt.Errorf("query %#q: %w", q, err)

		warn.Println(_f, err)

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
		warn.Println(_f, err)

		return nil, err
	}

	return locs, nil
}
