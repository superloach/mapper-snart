package main

import (
	"github.com/superloach/minori"

	"github.com/go-snart/bot"
	"github.com/go-snart/route"
)

var Log *minori.Logger

func Register(name string, b *bot.Bot) error {
	_f := "Register"
	Log = minori.GetLogger(name)

	Log.Info(_f, "forking registration")
	go register(name, b)
	Log.Info(_f, "forked registration")

	return nil
}

func register(name string, b *bot.Bot) error {
	_f := "register"
	Log.Info(_f, "registering routes")

	b.DB.Once(MapperDB)
	b.DB.Once(POITable)

	search := func(ctx *route.Ctx) error {
		return Search(b.DB, ctx)
	}

	b.Router.Add(
		&route.Route{
			Name:  "pois",
			Match: "pois?",
			Desc:  "Search for any POIs.",
			Cat:   name,
			Okay:  nil,
			Func:  search,
		},
		&route.Route{
			Name:  "gyms",
			Match: "g(yms?)?",
			Desc:  "Search for Pokemon Go gyms.",
			Cat:   name,
			Okay:  nil,
			Func:  search,
		},
		&route.Route{
			Name:  "stops",
			Match: "s(tops?)?",
			Desc:  "Search for Pokemon Go stops.",
			Cat:   name,
			Okay:  nil,
			Func:  search,
		},
	)

	Log.Info(_f, "registered routes")
	return nil
}
