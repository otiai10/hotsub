package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/otiai10/awsub/core"
	"github.com/otiai10/awsub/logs"
	"github.com/otiai10/awsub/parser"
	"github.com/otiai10/awsub/platform"
	"github.com/urfave/cli"
)

// action ...
// All the CLI context should be parsed and decoded on this layer,
// no deeper layer should NOT touch cli.
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

	root.Runtime.Image.Name = ctx.String("image")
	root.Runtime.Script.Path = ctx.String("script")

	jobs, err := parser.ParseFile(tasksfpath)
	if err != nil {
		return err
	}
	root.Jobs = jobs
	applog("Your tasks file is parsed and decoded to %d job(s) ✅", len(jobs))

	// {{{ TEMP
	root.JobLoggerer = new(logs.ColorfulLoggerFactory)
	// }}}

	shared, err := parser.ParseSharedData(ctx.StringSlice("shared"))
	if err != nil {
		return err
	}
	root.SharedData.Inputs = shared
	sdispec, err := platform.DefineSharedDataInstanceSpec(shared, ctx)
	if err != nil {
		return err
	}
	root.SharedData.Spec = sdispec

	spec, err := platform.DefineMachineSpec(ctx)
	if err != nil {
		return err
	}
	root.Machine.Spec = spec

	if err := platform.Get(ctx).Validate(); err != nil {
		return err
	}

	destroy := func() error {
		if !ctx.Bool("keep") {
			return root.Destroy()
		}
		return nil
	}

	if err := root.Create(); err != nil {
		destroy()
		return err
	}

	if err := root.Construct(); err != nil {
		destroy()
		return err
	}

	if err := root.Commit(nil); err != nil {
		destroy()
		return err
	}

	if err := destroy(); err != nil {
		return err
	}

	applog("All of your %d job(s) are completed 🎉", len(jobs))

	return nil
}

func applog(format string, v ...interface{}) {
	format = regexp.MustCompile("\n*$").ReplaceAllString(format, "\n")
	log.Printf("[COMMAND]\t"+format, v...)
}
