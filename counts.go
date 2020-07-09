package mapper

import "github.com/go-snart/snart/route"

func Counts(ctx *route.Ctx) error {
	rep := ctx.Reply()

	rep.Content = "stub"

	return rep.Send()
}
