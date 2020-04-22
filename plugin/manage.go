package main

import (
	"fmt"

	"github.com/go-snart/db"
	"github.com/go-snart/route"
)

func Manage(d *db.DB, ctx *route.Ctx) error {
	_f := "Manage"

	err := ctx.Flags.Parse()
	if err != nil {
		err = fmt.Errorf("flag parse: %w", err)
		Log.Error(_f, err)
		return err
	}

	rep := ctx.Reply()
	rep.Content = "manage stub"
	_, err = rep.Send()
	return err
}
