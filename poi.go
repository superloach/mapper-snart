package mapper

import "gopkg.in/rethinkdb/rethinkdb-go.v6/types"

type POI struct {
	ID    string `json:"id" rethinkdb:"id"`
	Name  string `json:"name" rethinkdb:"name"`
	Image string `json:"image" rethinkdb:"image"`
	Notes string `json:"notes" rethinkdb:"notes"`

	Loc types.Point `json:"loc" rethinkdb:"loc"`

	Ingr string `json:"ingr" rethinkdb:"ingr"`
	Pkmn string `json:"pkmn" rethinkdb:"pkmn"`
	Wzrd string `json:"wzrd" rethinkdb:"wzrd"`

	Alias []string `json:"alias" rethinkdb:"alias"`
}
