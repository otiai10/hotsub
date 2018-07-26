package core

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/otiai10/daap"
	"github.com/otiai10/ternary"
	"golang.org/x/sync/errgroup"
)

// Fetch downloads resources from cloud storage services,
// localize those URLs to local path,
// and, additionally, ensure the output directories exist.
func (job *Job) Fetch() error {

	eg := new(errgroup.Group)

	for _, input := range job.Parameters.Inputs {
		i := input
		eg.Go(func() error { return job.fetch(i) })
	}

	for _, output := range job.Parameters.Outputs {
		o := output
		eg.Go(func() error { return job.ensure(o) })
	}

	for _, include := range job.Parameters.Includes {
		i := include
		eg.Go(func() error { return job.upload(i) })
	}

	for _, env := range job.Parameters.Envs {
		job.addContainerEnv(env)
	}

	return eg.Wait()
}

// fetch
func (job *Job) fetch(input *Input) error {

	// Skip empty URL deliberately,
	// so that user defined script can switch its operation by itself.
	// TODO: This conditional skip should be moved to some validation steps in the future.
	if input.URL == "" {
		return nil
	}

	if err := input.Localize(HOTSUB_CONTAINERROOT); err != nil {
		return err
	}

	fetch := &daap.Execution{
		Inline:  "/scripts/download.sh",
		Env:     input.EnvForFetch(),
		Inspect: true,
	}

	ctx := context.Background()
	stream, err := job.Container.Routine.Exec(ctx, fetch)
	if err != nil {
		return err
	}

	for payload := range stream {
		job.Stdio(payload.Type, FETCH, payload.Text())
	}

	if fetch.ExitCode != 0 {
		return fmt.Errorf(
			"failed to download `%s` with status %d, please use --verbose option",
			input.URL, fetch.ExitCode,
		)
	}

	job.addContainerEnv(input.Env())

	return nil
}

// ensure the output directories exist on the workflow container.
func (job *Job) ensure(output *Output) error {
	// log.Println(job.Identity.Name, "ensure", output.URL)
	if err := output.Localize(HOTSUB_CONTAINERROOT); err != nil {
		return err
	}

	dir := ternary.If(output.Recursive).String(output.DeployedPath, filepath.Dir(output.DeployedPath))

	ensure := &daap.Execution{
		Inline:  fmt.Sprintf("mkdir -p %s", dir),
		Inspect: true,
	}

	ctx := context.Background()
	stream, err := job.Container.Routine.Exec(ctx, ensure)
	if err != nil {
		return err
	}

	for payload := range stream {
		job.Stdio(payload.Type, FETCH, payload.Text())
	}

	if ensure.ExitCode != 0 {
		return fmt.Errorf(
			"failed to download `%s` with status %d, please use --verbose option",
			output.URL, ensure.ExitCode,
		)
	}

	job.addContainerEnv(output.Env())
	return nil
}

// upload local files to the workflow container specified by "--include".
func (job *Job) upload(include *Include) error {

	ctx := context.Background()
	f, err := os.Open(include.LocalPath)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := job.Container.Workflow.Upload(ctx, f, HOTSUB_CONTAINERROOT); err != nil {
		return err
	}

	// This is just "Localize" specific for "Include".
	include.DeployedPath = filepath.ToSlash(filepath.Join(HOTSUB_CONTAINERROOT, filepath.Base(include.LocalPath)))

	// Skip anonymous includes
	if include.Resource.Name == "" {
		return nil
	}

	job.addContainerEnv(include.Env())
	return nil
}

// FIXME: Use channel to merge envs...
var lock = new(sync.Mutex)

func (job *Job) addContainerEnv(envs ...Env) {
	lock.Lock()
	defer lock.Unlock()
	job.Container.Envs = append(job.Container.Envs, envs...)
}
