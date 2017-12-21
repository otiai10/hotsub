// package main
// But, this module should be in a separated package in the future.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

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
		Driver:                      "amazonec2",
		AmazonEC2Region:             "ap-southeast-2",
		AmazonEC2IAMInstanceProfile: "testtest",
		Name: fmt.Sprintf("%s%02d", task.Prefix, task.Index),
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

	lifecycle := daap.NewContainer("awsub/lifecycle", daap.Args{
		Machine: &daap.MachineConfig{
			Host:     machine.Host(),
			CertPath: machine.CertPath(),
		},
		Mounts: []daap.Mount{
			daap.Volume("/tmp", "/tmp"),
		},
	})

	container := daap.NewContainer(h.Image, daap.Args{
		Machine: &daap.MachineConfig{
			Host:     machine.Host(),
			CertPath: machine.CertPath(),
		},
		Mounts: []daap.Mount{
			daap.Volume("/tmp", "/tmp"),
		},
	})

	ctx := context.Background()

	wg := new(sync.WaitGroup)
	wg.Add(2)
	go h.Warmup(ctx, lifecycle, job, wg)
	go h.Warmup(ctx, container, job, wg)
	wg.Wait()

	if job.Error != nil {
		return job
	}

	if err := h.Prepare(ctx, lifecycle, job); err != nil {
		return job.Errorf("failed to prepare input tasks: %v", err)
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

// Warmup ...
func (h *Handler) Warmup(ctx context.Context, c *daap.Container, job *Job, wg *sync.WaitGroup) {
	defer wg.Done()
	job.Logf("Pulling the image to the host: %v", c.Image)
	pull, err := c.PullImage(ctx)
	if err != nil {
		job.Errorf("failed to pull image: %v", err)
		return
	}
	for range pull {
		// fmt.Print(".")
	}
	job.Logf("The image pulled successfully: %v", c.Image)

	job.Logf("Creating a container on the host")
	if err = c.Create(ctx); err != nil {
		job.Errorf("failed to create container: %v: %v", c.Image, err)
		return
	}
	job.Logf("The container created successfully")

	job.Logf("Starting the container up")
	if err = c.Start(ctx); err != nil {
		job.Errorf("failed to start container: %v: %v", c.Image, err)
		return
	}
	job.Logf("The container started successfully: %v", c.Image)
}

// Prepare upload and locate files which are specified by the task onto the container.
// "Prepare" lifecycle is supposed to do
// 1) Download inputs files and directories specified by the task
// 2) Place those files on some specific location of the container.
// 3) Set the pairs of env variable and path to the file location to the task,
//    which is used by Execution.
func (h *Handler) Prepare(ctx context.Context, container *daap.Container, job *Job) error {

	task := job.Task
	for key, val := range task.Env {
		task.ContainerEnv = append(task.ContainerEnv, fmt.Sprintf("%s=%s", key, val))
	}
	envpairs := make(chan string)
	flag := 0 // flag!?
	for envname, url := range task.Inputs {
		flag++
		go h.prepareInput(ctx, container, envname, url, job, envpairs)
	}

	for envpair := range envpairs {
		flag--
		if envpair != "" {
			task.ContainerEnv = append(task.ContainerEnv, envpair)
		}
		if flag < 1 {
			close(envpairs)
		}
	}

	return nil
}

func (h *Handler) prepareInput(ctx context.Context, c *daap.Container, envname, url string, job *Job, result chan<- string) {
	stream, err := c.Exec(ctx, daap.Execution{
		Inline: "/lifecycle/download.sh",
		Env:    []string{fmt.Sprintf("%s=%s", "INPUT", url), fmt.Sprintf("%s=%s", "DIR", "/tmp")},
	})
	if err != nil {
		job.Errorf("failed to execute /lifecycle/download.sh: %v", err)
		result <- ""
		return
	}
	for payload := range stream {
		job.Logf("[PREPARE] &%d> %s", payload.Type, string(payload.Data))
	}
	result <- fmt.Sprintf("%s=%s", envname, filepath.Join("/tmp", filepath.Base(url)))
	return
}
