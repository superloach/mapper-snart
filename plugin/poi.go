package main

import r "gopkg.in/rethinkdb/rethinkdb-go.v6"

type POI struct {
	ID    string   `json:"id"    rethinkdb:"id"`
	Name  string   `json:"name"  rethinkdb:"name"`
	Image string   `json:"image" rethinkdb:"image"`
	Notes string   `json:"notes" rethinkdb:"notes"`
	Lat   float64  `json:"lat"   rethinkdb:"lat"`
	Lng   float64  `json:"lng"   rethinkdb:"lng"`
	Ingr  string   `json:"ingr"  rethinkdb:"ingr"`
	Pkmn  string   `json:"pkmn"  rethinkdb:"pkmn"`
	Wzrd  string   `json:"wzrd"  rethinkdb:"wzrd"`
	Alias []string `json:"alias" rethinkdb:"alias"`
}

var POITable = r.DB("poi").TableCreate(
	"poi",
	r.TableCreateOpts{
		PrimaryKey: "id",
	},
)
