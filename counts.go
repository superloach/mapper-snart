package mapper

import "github.com/go-snart/snart/route"

// Counts returns the (bounded?) POI totals for various groups.
func Counts(ctx *route.Ctx) error {
	rep := ctx.Reply()

	rep.Content = "stub"

	return rep.Send()
}
