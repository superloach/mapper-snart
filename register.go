package mapper

import (
	"github.com/superloach/minori"

	"github.com/go-snart/snart/bot"
	"github.com/go-snart/snart/route"
)

var Log = minori.GetLogger("mapper")

func init() {
	bot.Register("mapper", Register)
}

func Register(b *bot.Bot) error {
	_f := "Register"
	Log.Info(_f, "registering")

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
			Cat:   "mapper",
			Okay:  nil,
			Func:  search,
		},
		&route.Route{
			Name:  "gyms",
			Match: "g(yms?)?",
			Desc:  "Search for Pokemon Go gyms.",
			Cat:   "mapper",
			Okay:  nil,
			Func:  search,
		},
		&route.Route{
			Name:  "stops",
			Match: "s(tops?)?",
			Desc:  "Search for Pokemon Go stops.",
			Cat:   "mapper",
			Okay:  nil,
			Func:  search,
		},
	)

	Log.Info(_f, "registered")

	return nil
}
