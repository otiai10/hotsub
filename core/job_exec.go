package core

import (
	"context"
	"fmt"

	"github.com/otiai10/daap"
)

// Exec executes user-defined script inside the workflow container.
func (job *Job) Exec() error {

	ctx := context.Background()

	env := []string{fmt.Sprintf("%s=%s", "HOTSUB_ROOT", HOTSUB_CONTAINERROOT)}
	for _, e := range job.Container.Envs {
		env = append(env, e.Pair())
	}

	workflow := &daap.Execution{
		Script:  job.Container.Script.Path,
		Env:     env,
		Inspect: true,
	}

	stream, err := job.Container.Workflow.Exec(ctx, workflow)
	if err != nil {
		return err
	}

	for payload := range stream {
		job.Stdio(payload.Type, EXECUTE, string(payload.Data))
	}

	if workflow.ExitCode != 0 {
		return fmt.Errorf("script exited with %d, check verbose log", workflow.ExitCode)
	}

	return nil
}
