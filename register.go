package mapper

import (
	"github.com/superloach/minori"

	dg "github.com/bwmarrin/discordgo"
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

	b.AddGamer(bot.GamerText(
		"Pokémon GO",
		dg.GameTypeGame,
	))
	b.AddGamer(GamerCounts(
		"%.f Gyms | %.f PokéStops",
		map[string]interface{}{"pkmn": "G"},
		map[string]interface{}{"pkmn": "S"},
	))

	b.AddGamer(bot.GamerText(
		"Ingress",
		dg.GameTypeGame,
	))
	b.AddGamer(GamerCounts(
		"%.f POIs",
		map[string]interface{}{},
	))

	search := func(ctx *route.Ctx) error {
		return Search(b.DB, ctx, b.Admin(ctx))
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
		&route.Route{
			Name:  "qtest",
			Match: "qtest",
			Desc:  "test command for queryer",
			Cat:   "mapper",
			Okay:  nil,
			Func:  b.DB.Queryer(Qtest),
		},
	)

	Log.Info(_f, "registered")

	return nil
}
