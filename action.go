package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

// action ...
func action(ctx *cli.Context) error {

	if ctx.NumFlags() == 0 {
		return ctx.App.Command("help").Run(ctx)
	}

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

	handler, err := NewHandler(ctx)
	if err != nil {
		return err
	}

	errored := []*Job{}
	for report := range handler.HandleBunch(tasks) {
		if report.Error != nil {
			errored = append(errored, report)
		}
	}

	if len(errored) == 0 {
		fmt.Printf("All %d tasks completed successfully!\n", len(tasks))
	} else {
		fmt.Printf("%d task(s) failed with errors.\n", len(errored))
		for _, job := range errored {
			fmt.Printf("%s: %v\n", job.Instance.Name, job.Error)
		}
	}

	return nil
}
