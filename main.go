package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Flags = flags
	app.Action = action
	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}
