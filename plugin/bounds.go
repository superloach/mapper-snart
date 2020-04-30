package main

import (
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
	"gopkg.in/rethinkdb/rethinkdb-go.v6/types"
)

type Bounds struct {
	ID   string      `json:"id"   rethinkdb:"id"`
	Poly types.Lines `json:"poly" rethinkdb:"poly"`
}

var BoundsTable = r.DB("poi").TableCreate(
	"bounds",
	r.TableCreateOpts{
		PrimaryKey: "id",
	},
)
