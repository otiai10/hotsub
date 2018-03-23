// package main
// But, this module should be in a separated package in the future.
package main

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/otiai10/daap"
	"github.com/otiai10/dkmachine/v0/dkmachine"
	"github.com/urfave/cli"
)

// Handler ...
type Handler struct {
	Image   string
	Script  string
	Verbose bool
	ctx     *cli.Context
}

// NewHandler ...
func NewHandler(ctx *cli.Context) (*Handler, error) {
	h := &Handler{ctx: ctx}

	h.Image = ctx.String("image")
	if h.Image == "" {
		return nil, fmt.Errorf("`--image` is required but not specified")
	}

	script := ctx.String("script")
	if script == "" {
		return nil, fmt.Errorf("`--script` is required (for now)")
	}
	script, err := filepath.Abs(script)
	if err != nil {
		return nil, err
	}
	h.Script = script

	h.Verbose = ctx.Bool("verbose")

	return h, nil
}

// HandleBunch ...
func (h *Handler) HandleBunch(tasks []*Task) <-chan *Job {

	results := make(chan *Job)
	wg := new(sync.WaitGroup)
	wg.Add(len(tasks))

	for _, task := range tasks {
		time.Sleep(time.Millisecond * time.Duration(h.ctx.Int("debug-sleep")))
		go func(t *Task) {
			results <- h.Handle(t)
			wg.Done()
		}(task)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}

// Handle ...
func (h *Handler) Handle(task *Task) *Job {

	job := &Job{Task: task}

	instance, err := h.generateMachineOption(task)
	if err != nil {
		return job.Errorf("failed to specify machine config: %v", err)
	}
	job.Instance = instance

	if h.Verbose {
		job.Logger = NewLogger(fmt.Sprintf("[%s]", job.Instance.Name), task.Index)
	}

	job.Logf("Creating docker machine on %v", strings.ToUpper(instance.Driver))
	machine, err := dkmachine.Create(job.Instance)
	if machine != nil && !h.ctx.Bool("keep") {
		defer func() {
			job.Logf("Deleted machine: %v", machine.Remove())
		}()
	}
	if err != nil {
		return job.Errorf("failed to create machine: %v", err)
	}
	job.Logf("The machine created successfully")

	lifecycle := daap.NewContainer("awsub/lifecycle", machine)
	maincontainer := daap.NewContainer(h.Image, machine)

	ctx := context.Background()

	wg := new(sync.WaitGroup)
	wg.Add(2)
	go h.Warmup(ctx, lifecycle, job, wg)
	go h.Warmup(ctx, maincontainer, job, wg)
	wg.Wait()

	if job.Error != nil {
		return job
	}

	if err := h.Prepare(ctx, lifecycle, job); err != nil {
		return job.Errorf("failed to prepare input tasks: %v", err)
	}

	// Add super variables
	task.ContainerEnv = append(task.ContainerEnv, fmt.Sprintf("%s=%s", "AWSUB_ROOT", AWSUBROOT))

	execution := &daap.Execution{
		Script:  h.Script,
		Env:     task.ContainerEnv,
		Inspect: true,
	}

	job.Logf("Sending command queue to the container")
	stream, err := maincontainer.Exec(ctx, execution)
	if err != nil {
		return job.Errorf("failed to exec script in the container: %v", err)
	}
	job.Logf("The command queued successfully")

	for payload := range stream {
		job.Logf("&%d> %s", payload.Type, string(payload.Data))
	}
	if execution.ExitCode != 0 {
		return job.Errorf("your script exited with code %d, please check verbose log", execution.ExitCode)
	}
	job.Logf("The command finished completely")

	if err := h.Finalize(ctx, lifecycle, job); err != nil {
		return job.Errorf("failed to finalize output of task: %v", err)
	}
	job.Logf("Successfully uploaded output files to your bucket")

	job.Logf("Everything completed. Good work :)")

	return job
}
