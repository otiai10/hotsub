package main

import (
	"context"
	"fmt"
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
