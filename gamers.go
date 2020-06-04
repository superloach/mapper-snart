package mapper

import (
	"fmt"

	dg "github.com/bwmarrin/discordgo"
	"github.com/go-snart/snart/bot"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

func GamerText(text string) bot.Gamer {
	return func(b *bot.Bot) (*dg.Game, error) {
		return &dg.Game{
			Name: text,
			Type: dg.GameTypeGame,
		}, nil
	}
}

func GamerCount(filt interface{}, lbl string) bot.Gamer {
	return func(b *bot.Bot) (*dg.Game, error) {
		q := r.DB("mapper").Table("poi").Filter(filt).Count()
		count := make([]int, 0)
		err := q.ReadAll(&count, b.DB)
		if err != nil {
			return nil, err
		}

		return &dg.Game{
			Name: fmt.Sprintf(lbl, count[0]),
			Type: dg.GameTypeWatching,
		}, nil
	}
}
