package main

import (
	"github.com/otiai10/hotsub/application"
	"github.com/otiai10/hotsub/flags"
	"github.com/urfave/cli"
)

var commands = []cli.Command{

	// run
	cli.Command{
		Name:        "run",
		Description: "Run your jobs on cloud with specified input files and any parameters",
		Usage:       "Run your jobs on cloud with specified input files and any parameters",
		Action: func(ctx *cli.Context) error {
			if ctx.NumFlags() == 0 {
				return ctx.App.Command("help").Run(ctx)
			}
			return application.Run(ctx)
		},
		Flags: flags.Index,
	},
}
