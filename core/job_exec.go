package core

import (
	"context"
	"fmt"

	"github.com/otiai10/daap"
)

// Exec executes user-defined script inside the workflow container.
func (job *Job) Exec() error {

	ctx := context.Background()

	workflow := job.toWorkflowExecution()

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

func (job *Job) toWorkflowExecution() *daap.Execution {

	env := []string{fmt.Sprintf("%s=%s", "HOTSUB_ROOT", HotsubContainerRoot)}
	for _, e := range job.Container.Envs {
		env = append(env, e.Pair())
	}

	workflow := &daap.Execution{Inspect: true, Env: env}

	switch job.Type {
	case CommonWorkflowLanguageJob:
		workflow.Inline = "cwltool ${CWL_FILE} ${CWL_PARAM_FILE}"
	default:
		workflow.Script = job.Container.Script.Path
	}

	return workflow
}
