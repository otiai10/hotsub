package application

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/otiai10/hotsub/core"
	"github.com/otiai10/hotsub/logs"
	"github.com/otiai10/hotsub/params"
	"github.com/otiai10/hotsub/parser"
	"github.com/otiai10/hotsub/platform"
)

func generateJobsFromContext(ctx params.Context) (string, []*core.Job, error) {

	// {{{ FIXME: ugly
	if cwlfile := ctx.String("cwl"); cwlfile != "" {
		name := filepath.Base(cwlfile)
		jobs := []*core.Job{}
		for index, jobparam := range ctx.StringSlice("cwl-job") {
			job := core.NewJob(index, name)
			job.Parameters.Includes = core.Includes{
				{LocalPath: cwlfile, Resource: core.Resource{Name: "CWL_FILE"}},
				{LocalPath: jobparam, Resource: core.Resource{Name: "CWL_JOB_FILE"}},
			}
			job.Type = core.CommonWorkflowLanguageJob
			jobs = append(jobs, job)
		}
		return name, jobs, nil
	}

	if wdlfile := ctx.String("wdl"); wdlfile != "" {
		name := filepath.Base(wdlfile)
		jobs := []*core.Job{}
		for index, jobparam := range ctx.StringSlice("wdl-job") {
			job := core.NewJob(index, name)
			job.Parameters.Includes = core.Includes{
				{LocalPath: wdlfile, Resource: core.Resource{Name: "WDL_FILE"}},
				{LocalPath: jobparam, Resource: core.Resource{Name: "WDL_JOB_FILE"}},
			}
			job.Type = core.WorkflowDescriptionLanguageJob
			jobs = append(jobs, job)
		}
		return name, jobs, nil
	}
	// }}}

	// Get tasks file path from context.
	tasksfpath := ctx.String("tasks")
	f, err := os.Open(tasksfpath)
	if err != nil {
		return "", nil, fmt.Errorf("failed to open tasks file `%s`: %v", tasksfpath, err)
	}
	defer f.Close()

	// Create jobs model from tasks file.
	jobs, err := parser.ParseFile(tasksfpath)
	if err != nil {
		return "", nil, err
	}

	// Decide workflow name by tasks file name.
	name := filepath.Base(tasksfpath)

	return name, jobs, nil
}

// Run ...
func Run(ctx params.Context) error {

	name, jobs, err := generateJobsFromContext(ctx)
	if err != nil {
		return err
	}

	// {{{ Add included files to each job.
	includes := parser.ParseIncludes(ctx.StringSlice("include"))
	for _, job := range jobs {
		job.Parameters.Includes = append(job.Parameters.Includes, includes...)
	}
	// }}}

	root := core.RootComponentTemplate(name)
	root.Jobs = jobs

	root.Runtime.Image.Name = ctx.String("image")

	// {{{ FIXME: ugly
	if len(jobs) != 0 {
		switch jobs[0].Type {
		case core.CommonWorkflowLanguageJob:
			root.Runtime.Image.Name = "hotsub/c4cwl"
		case core.WorkflowDescriptionLanguageJob:
			root.Runtime.Image.Name = "hotsub/c4wdl"
		}
	}
	// }}}

	root.Runtime.Script.Path = ctx.String("script")
	root.Concurrency = ctx.Int64("concurrency")

	applog("Your tasks file is parsed and decoded to %d job(s) ✅", len(jobs))

	// {{{ Define Log Location
	// root.JobLoggerFactory = new(logs.ColorfulLoggerFactory)
	dir := ctx.String("log-dir")
	if dir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		dir = filepath.Join(cwd, "log", time.Now().Format("20060102_150405"))
	}
	factory := &logs.IntegratedLoggerFactory{File: &logs.FileLoggerFactory{Dir: dir}}
	if ctx.Bool("verbose") {
		factory.Verbose = new(logs.ColorfulLoggerFactory)
	}
	root.JobLoggerFactory = factory
	log.Printf("[COMMAND]\tSee logs here -> %s\n", dir)
	// }}}

	if err := platform.Get(ctx).Validate(ctx); err != nil {
		return err
	}

	commonEnv, err := parser.ParseEnv(ctx.StringSlice("env"))
	if err != nil {
		return err
	}
	root.CommonParameters.Envs = commonEnv

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

	if err := root.Prepare(); err != nil {
		return err
	}

	destroy := func() error {
		if !ctx.Bool("keep") {
			return root.Destroy()
		}
		return nil
	}

	rootctx := context.Background()
	if err := root.Run(rootctx); err != nil {
		destroy()
		return err
	}

	if err := destroy(); err != nil {
		return err
	}

	applog("All of your %d job(s) are completed 🎉", len(jobs))

	return nil
}
