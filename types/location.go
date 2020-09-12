package types

import (
	"fmt"

	"github.com/dewski/spatial"

	"github.com/go-snart/snart/db"
)

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
func GetLocations(d *db.DB) ([]*Location, error) {
	return nil, fmt.Errorf("stub")
}
