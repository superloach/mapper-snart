package mapper

import (
	"fmt"

	dg "github.com/bwmarrin/discordgo"
	"github.com/go-snart/snart/bot"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

// GamerCounts creates a bot.Gamer which displays Location counts.
func GamerCounts(lbl string, filts ...interface{}) bot.Gamer {
	_f := "GamerCounts"

	return func(b *bot.Bot) (*dg.Game, error) {
		loct := LocationTable.Build(b.DB)

		filtqs := make([]r.Term, len(filts))
		for i, filt := range filts {
			filtqs[i] = loct.Filter(filt).Count()
		}

		counts := make([]interface{}, len(filts))

		for i, filtq := range filtqs {
			tmp := make([]interface{}, 0)
			err := filtq.ReadAll(&tmp, b.DB)

			if err != nil {
				err = fmt.Errorf("readall tmp: %w", err)
				Log.Error(_f, err)

				return nil, err
			}

			counts[i] = tmp[0]
		}

		return &dg.Game{
			Name: fmt.Sprintf(lbl, counts...),
			Type: dg.GameTypeWatching,
		}, nil
	}
}
