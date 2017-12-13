package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

// action ...
func action(ctx *cli.Context) error {

	tasksfpath := ctx.String("tasks")
	f, err := os.Open(tasksfpath)
	if err != nil {
		return fmt.Errorf("failed to open tasks file `%s`: %v", tasksfpath, err)
	}
	defer f.Close()

	tasks, err := parseTasksFromFile(f)
	if err != nil {
		return err
	}

	for report := range handleBunch(tasks) {
		fmt.Printf("%+v\n", report)
	}

	return nil
}
