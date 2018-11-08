package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	cli.VersionFlag = cli.BoolFlag{
		Name:  "version,V",
		Usage: "print the version",
	}
	app.Commands = commands
	app.Version = version
	app.Usage = "command line to run batch computing on AWS and GCP with the same interface"
	app.Description = "Open-source command-line tool to run batch computing tasks and workflows on backend services such as Amazon Web Services."
	// app.Flags = flags.Index
	if err := app.Run(os.Args); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
