package main

import (
	"context"
	"fmt"
	"net/url"
	"path/filepath"

	"github.com/otiai10/daap"
)

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
	for envname, rawurl := range task.Inputs {
		flag++
		go h.prepareInput(ctx, container, envname, rawurl, job, envpairs)
	}

	for envname, rawurl := range task.InputRecursive {
		flag++
		go h.prepareInputRecursive(ctx, container, envname, rawurl, job, envpairs)
	}

	for envname, rawurl := range task.OutputRecursive {
		flag++
		go h.prepareOutputDirectory(ctx, container, envname, rawurl, job, envpairs)
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

func (h *Handler) prepareInput(ctx context.Context, c *daap.Container, envname, rawurl string, job *Job, result chan<- string) {
	stream, err := c.Exec(ctx, daap.Execution{
		Inline: "/lifecycle/download.sh",
		Env:    []string{fmt.Sprintf("%s=%s", "INPUT", rawurl), fmt.Sprintf("%s=%s", "DIR", "/tmp")},
	})
	if err != nil {
		job.Errorf("failed to execute /lifecycle/download.sh: %v", err)
		result <- ""
		return
	}
	for payload := range stream {
		job.Logf("[PREPARE] &%d> %s", payload.Type, string(payload.Data))
	}
	result <- fmt.Sprintf("%s=%s", envname, filepath.Join("/tmp", filepath.Base(rawurl)))
	return
}

func (h *Handler) prepareInputRecursive(ctx context.Context, c *daap.Container, envname, rawurl string, job *Job, result chan<- string) {
	stream, err := c.Exec(ctx, daap.Execution{
		Inline: "/lifecycle/download.sh",
		Env:    []string{fmt.Sprintf("%s=%s", "INPUT_RECURSIVE", rawurl), fmt.Sprintf("%s=%s", "DIR", "/tmp")},
	})
	if err != nil {
		job.Errorf("failed to execute /lifecycle/download.sh: %v", err)
		result <- ""
		return
	}
	for payload := range stream {
		job.Logf("[PREPARE] &%d> %s", payload.Type, string(payload.Data))
	}
	result <- fmt.Sprintf("%s=%s", envname, filepath.Join("/tmp", filepath.Base(rawurl)))
	return
}

func (h *Handler) prepareOutputDirectory(ctx context.Context, c *daap.Container, envname, rawurl string, job *Job, result chan<- string) {
	u, err := url.Parse(rawurl)
	if err != nil {
		result <- ""
		return
	}
	outdir := filepath.Join("/tmp", u.Path)
	stream, err := c.Exec(ctx, daap.Execution{
		Inline: fmt.Sprintf("mkdir -p %v", outdir),
	})
	if err != nil {
		job.Errorf("failed to execute make output directory: %v", err)
		result <- ""
		return
	}
	for payload := range stream {
		job.Logf("[PREPARE] &%d> %s", payload.Type, string(payload.Data))
	}
	result <- fmt.Sprintf("%s=%s", envname, outdir)
}
