package flags

import "github.com/urfave/cli"

// Template is a flag list for "template" subcommand.
var Template = []cli.Flag{
	cli.StringFlag{
		Name:  "name,n",
		Usage: "Template project name",
		Value: "helloworld",
	},
	cli.StringFlag{
		Name:  "dir,d",
		Usage: "Directory to create this template",
		Value: ".",
	},
}
