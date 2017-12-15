// package main
// But, this module should be in a separated package in the future.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/otiai10/daap"
	"github.com/otiai10/dkmachine/v0/dkmachine"
	"github.com/urfave/cli"
)

// Handler ...
type Handler struct {
	Image   string
	Script  string
	Verbose bool
}

// NewHandler ...
func NewHandler(ctx *cli.Context) (*Handler, error) {
	h := &Handler{}
	img := ctx.String("image")
	if img == "" {
		return nil, fmt.Errorf("`--image` is required but not specified")
	}
	h.Image = img
	script := ctx.String("script")
	if script == "" {
		return nil, fmt.Errorf("`--script` is required (for now)")
	}
	script, err := filepath.Abs(script)
	if err != nil {
		return nil, err
	}
	h.Script = script

	// {{{ FIXME: debug
	h.Verbose = true
	// }}}

	return h, nil
}

// HandleBunch ...
func (h *Handler) HandleBunch(tasks []*Task) <-chan *Job {

	results := make(chan *Job)
	done := make(chan bool)
	count := 0

	for _, task := range tasks {
		go func(t *Task) {
			results <- h.Handle(t)
			done <- true
		}(task)
	}

	go func() {
		defer close(done)
		defer close(results)
		for {
			<-done
			count++
			if count >= len(tasks) {
				return
			}
		}
	}()

	return results
}

// Handle ...
func (h *Handler) Handle(task *Task) *Job {

	job := &Job{Task: task}

	// {{{ TODO: Refactor and Slim up
	job.Instance = &dkmachine.CreateOptions{
		Driver:          "amazonec2",
		AmazonEC2Region: "ap-southeast-2",
		Name:            fmt.Sprintf("%s%03d", task.Prefix, task.Index),
	}

	// {{{ FIXME: debug
	job.Logger = log.New(os.Stdout, fmt.Sprintf("[%s] ", job.Instance.Name), 0)
	// }}}

	job.Logf("Creating docker machine")
	machine, err := dkmachine.Create(job.Instance)
	if err != nil {
		return job.Errorf("failed to create machine: %v", err)
	}
	job.Logf("The machine created successfully")
	defer func() {
		job.Logf("Deleting docker machine")
		machine.Remove()
		job.Logf("The machine deleted successfully")
	}()

	container := daap.NewContainer(h.Image, daap.Args{
		Machine: &daap.MachineConfig{
			Host:     machine.Host(),
			CertPath: machine.CertPath(),
		},
		Env: []string{},
	})

	ctx := context.Background()

	job.Logf("Pulling the image to the host: %v", container.Image)
	pull, err := container.PullImage(ctx)
	if err != nil {
		return job.Errorf("failed to pull image: %v", err)
	}
	for range pull {
		// fmt.Print(".")
	}
	job.Logf("The image pulled successfully")

	job.Logf("Creating a container on the host")
	if err = container.Create(ctx); err != nil {
		return job.Errorf("failed to create container: %v", err)
	}
	job.Logf("The container created successfully")

	job.Logf("Starting the container up")
	if err = container.Start(ctx); err != nil {
		return job.Errorf("failed to start container: %v", err)
	}
	job.Logf("The container started successfully")

	if err := h.Prepare(container, job); err != nil {
		return job.Errorf("failed to prepare files specified by task: %v", err)
	}

	execution := daap.Execution{
		Script: h.Script,
		Env:    task.ContainerEnv,
	}

	job.Logf("Sending command queue to the container")
	stream, err := container.Exec(ctx, execution)
	if err != nil {
		return job.Errorf("failed to exec script in the container: %v", err)
	}
	job.Logf("The command queued successfully")

	for payload := range stream {
		job.Logf("&%d> %s", payload.Type, string(payload.Data))
	}
	job.Logf("The command finished completely")
	// }}}

	return job
}

// Prepare upload and locate files which are specified by the task onto the container.
// "Prepare" lifecycle is supposed to do
// 1) Download inputs files and directories specified by the task
// 2) Place those files on some specific location of the container.
// 3) Set the pairs of env variable and path to the file location to the task,
//    which is used by Execution.
func (h *Handler) Prepare(container *daap.Container, job *Job) error {
	task := job.Task
	for key, val := range task.Env {
		task.ContainerEnv = append(task.ContainerEnv, fmt.Sprintf("%s=%s", key, val))
	}
	for key, val := range task.Inputs {
		fmt.Println(key, val)
	}
	for key, val := range task.InputRecursive {
		fmt.Println(key, val)
	}
	return nil
}
