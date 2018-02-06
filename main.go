package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

const (
	version = "0.0.6"
)

func main() {
	app := cli.NewApp()
	cli.VersionFlag = cli.BoolFlag{
		Name:  "version,V",
		Usage: "print the version",
	}
	app.Commands = commands
	app.Version = version
	app.Usage = "command line to run batch computing on AWS"
	app.Description = "Open-source command-line tool to run batch computing tasks and workflows on backend services such as Amazon Web Service."
	app.Flags = flags
	app.Action = action
	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
