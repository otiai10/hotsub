package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/otiai10/awsub/core"
	"github.com/otiai10/awsub/parser"
	"github.com/otiai10/awsub/platform"
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

	name := filepath.Base(tasksfpath)
	root := core.RootComponentTemplate(name)

	jobs, err := parser.ParseFile(tasksfpath)
	if err != nil {
		return err
	}
	root.Jobs = jobs

	spec, err := platform.DefineMachineSpec(ctx)
	if err != nil {
		return err
	}
	root.Machine.Spec = spec

	if err := platform.Get(ctx).Validate(); err != nil {
		return err
	}

	if err := root.Commit(nil); err != nil {
		return err
	}

	// handler, err := NewHandler(ctx)
	// if err != nil {
	// 	return err
	// }
	// errored := []*Job{}
	//
	// for report := range handler.HandleBunch(tasks) {
	// 	if report.Error != nil {
	// 		errored = append(errored, report)
	// 	}
	// }
	//
	// if len(errored) == 0 {
	// 	fmt.Printf("All %d tasks completed successfully!\n", len(tasks))
	// } else {
	// 	for _, job := range errored {
	// 		fmt.Printf("%s: %v\n", job.Instance.Name, job.Error)
	// 	}
	// 	return fmt.Errorf("%d task(s) failed with errors", len(errored))
	// }

	return nil
}
