package mapper

import (
	"fmt"

	dg "github.com/bwmarrin/discordgo"
	"github.com/go-snart/snart/bot"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

func GamerPokemonGO(b *bot.Bot) (*dg.Game, error) {
	return &dg.Game{
		Name: "Pokémon GO",
		Type: dg.GameTypeGame,
	}, nil
}

func GamerGyms(b *bot.Bot) (*dg.Game, error) {
	q := r.DB("mapper").Table("poi").Filter(
		map[string]interface{}{
			"pkmn": "G",
		},
	).Count()

	count := 0
	err := q.ReadAll(&count, b.DB)
	if err != nil {
		return nil, err
	}

	return &dg.Game{
		Name: fmt.Sprintf("%d gyms", count),
		Type: dg.GameTypeWatching,
	}, nil
}
