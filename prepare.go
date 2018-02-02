package main

import (
	"context"
	"fmt"
	"net/url"
	"path/filepath"

	"github.com/otiai10/daap"
)

const (
	// AWSUBROOT ...
	AWSUBROOT = "/tmp"
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

	if flag == 0 {
		close(envpairs)
		return job.Error
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

	return job.Error
}

// prepareInput prepares input files to the worker container and convert it
// to envirionment variables for the execution of user script.
// rawurl is a raw URL with s3:// or gs:// scheme.
// Every input file will be located on a path which is generated just by
// replacing scheme and bucket name of rawurl with "workdir".
// e.g.
//  s3://hgc-otiai10-test/foo/bar.txt -> ${workdir}/foo/bar.txt
func (h *Handler) prepareInput(ctx context.Context, c *daap.Container, envname, rawurl string, job *Job, result chan<- string) {

	if rawurl == "" {
		result <- ""
		return
	}

	u, err := url.Parse(rawurl)
	if err != nil {
		job.Errorf("failed to parse the url of given input: %s: %v", rawurl, err)
		result <- ""
		return
	}

	dir := filepath.Join(AWSUBROOT, u.Hostname(), filepath.Dir(u.Path))

	execution := &daap.Execution{
		Inline:  "/lifecycle/download.sh",
		Env:     []string{fmt.Sprintf("%s=%s", "INPUT", u.String()), fmt.Sprintf("%s=%s", "DIR", dir)},
		Inspect: true,
	}
	stream, err := c.Exec(ctx, execution)
	if err != nil {
		job.Errorf("failed to execute /lifecycle/download.sh: %v", err)
		result <- ""
		return
	}
	for payload := range stream {
		job.Logf("[PREPARE] &%d> %s", payload.Type, string(payload.Data))
	}

	if execution.ExitCode != 0 {
		job.Errorf("failed to download input file `%s` with status code %d, please check output with --verbose option", rawurl, execution.ExitCode)
		result <- ""
		return
	}

	result <- fmt.Sprintf("%s=%s", envname, filepath.Join(dir, filepath.Base(rawurl)))
	return
}

// prepareInput prepares input directories to the worker container and convert it
// to envirionment variables for the execution of user script.
// rawurl is a raw URL with s3:// or gs:// scheme.
// Every input file will be located on a path which is generated just by
// replacing scheme and bucket name of rawurl with "workdir".
// e.g.
//  s3://hgc-otiai10-test/foo/baz -> ${workdir}/foo/baz
func (h *Handler) prepareInputRecursive(ctx context.Context, c *daap.Container, envname, rawurl string, job *Job, result chan<- string) {

	if rawurl == "" {
		result <- ""
		return
	}

	u, err := url.Parse(rawurl)
	if err != nil {
		job.Errorf("failed to parse the url of given input: %s: %v", rawurl, err)
		result <- ""
		return
	}

	dir := filepath.Join(AWSUBROOT, u.Hostname(), filepath.Dir(u.Path))
	execution := &daap.Execution{
		Inline:  "/lifecycle/download.sh",
		Env:     []string{fmt.Sprintf("%s=%s", "INPUT_RECURSIVE", rawurl), fmt.Sprintf("%s=%s", "DIR", dir)},
		Inspect: true,
	}
	stream, err := c.Exec(ctx, execution)
	if err != nil {
		job.Errorf("failed to execute /lifecycle/download.sh: %v", err)
		result <- ""
		return
	}
	for payload := range stream {
		job.Logf("[PREPARE] &%d> %s", payload.Type, string(payload.Data))
	}

	if execution.ExitCode != 0 {
		job.Errorf("failed to download input file `%s` with status code %d, please check output with --verbose option", rawurl, execution.ExitCode)
		result <- ""
		return
	}

	result <- fmt.Sprintf("%s=%s", envname, filepath.Join(dir, filepath.Base(rawurl)))
	return
}

func (h *Handler) prepareOutputDirectory(ctx context.Context, c *daap.Container, envname, rawurl string, job *Job, result chan<- string) {
	u, err := url.Parse(rawurl)
	if err != nil {
		result <- ""
		return
	}
	outdir := filepath.Join(AWSUBROOT, u.Hostname(), u.Path)
	execution := &daap.Execution{
		Inline: fmt.Sprintf("mkdir -p %v", outdir),
	}
	stream, err := c.Exec(ctx, execution)
	if err != nil {
		job.Errorf("failed to execute make output directory: %v", err)
		result <- ""
		return
	}
	for payload := range stream {
		job.Logf("[PREPARE] &%d> %s", payload.Type, string(payload.Data))
	}

	if execution.ExitCode != 0 {
		job.Errorf("failed to create output directory `%s` with status code %d, please check output with --verbose option", rawurl, execution.ExitCode)
		result <- ""
		return
	}

	result <- fmt.Sprintf("%s=%s", envname, outdir)
}
