package main

import (
	"github.com/superloach/minori"

	"github.com/go-snart/snart/lib/bot"
	"github.com/go-snart/snart/lib/errs"
	"github.com/go-snart/snart/lib/plugin"
	"github.com/go-snart/snart/lib/route"
)

var Log *minori.Logger

func Register(name string, b *bot.Bot) error {
	_f := "Register"

	Log = plugin.Log.GetLogger(name)
	Log.Info(_f, "forking registration")

	go func() {
		b.DB.Easy(POIDB)
		b.DB.Easy(POITable)

		err := routes(name, b)
		if err != nil {
			errs.Wrap(&err, `routes(%#v, b)`, name)
			Log.Warn(_f, err)
			return
		}
	}()

	Log.Info(_f, "forked registration")
	return nil
}

func routes(name string, b *bot.Bot) error {
	_f := "routes"
	Log.Info(_f, "registering routes")

	poi := func(ctx *route.Ctx) error {
		err := Poi(b.DB, ctx)
		if err != nil {
			errs.Wrap(&err, `_poi(d, ctx)`)
			Log.Error(_f, err)
			return err
		}

		return nil
	}

	b.AddRoute(
		&route.Route{
			Name:  "pois",
			Match: "p(ois?)?",
			Desc:  "Search for any POIs. (Alias: `poi`, `p`)",
			Cat:   name,
			Okay:  nil,
			Func:  poi,
		},
		&route.Route{
			Name:  "gyms",
			Match: "g(yms?)?",
			Desc:  "Search for Pokemon Go gyms. (Alias: `gym`, `g`)",
			Cat:   name,
			Okay:  nil,
			Func:  poi,
		},
		&route.Route{
			Name:  "stops",
			Match: "s(tops?)?",
			Desc:  "Search for Pokemon Go stops. (Alias: `stop`, `s`)",
			Cat:   name,
			Okay:  nil,
			Func:  poi,
		},
	)

	Log.Info(_f, "registered routes")
	return nil
}
