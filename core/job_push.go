package core

import (
	"context"
	"fmt"

	"github.com/otiai10/daap"
	"golang.org/x/sync/errgroup"
)

// Push uploads result files to cloud storage services according to
// specified "output" URLs of tasks file.
func (job *Job) Push() error {

	eg := new(errgroup.Group)
	for _, output := range job.Parameters.Outputs {
		o := output
		eg.Go(func() error { return job.push(o) })
	}

	return eg.Wait()
}

func (job *Job) push(output *Output) error {
	ctx := context.Background()

	push := &daap.Execution{
		Inline: "/lifecycle/upload.sh",
		Env: []string{
			fmt.Sprintf("%s=%s", "SOURCE", output.LocalPath),
			fmt.Sprintf("%s=%s", "DEST", output.URL),
		},
		Inspect: true,
	}

	stream, err := job.Container.Routine.Exec(ctx, push)
	if err != nil {
		return err
	}

	for payload := range stream {
		fmt.Printf("&%d> %s\n", payload.Type, payload.Text())
	}

	if push.ExitCode != 0 {
		return fmt.Errorf("failed to upload `%s` with %d", output.URL, push.ExitCode)
	}

	return nil
}
