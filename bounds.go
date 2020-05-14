package main

import "gopkg.in/rethinkdb/rethinkdb-go.v6/types"

type Bounds struct {
	ID   string      `json:"id"   rethinkdb:"id"`
	Poly types.Lines `json:"poly" rethinkdb:"poly"`
}
