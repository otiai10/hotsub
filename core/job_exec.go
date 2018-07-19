package core

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/otiai10/daap"
)

// Exec executes user-defined script inside the workflow container.
func (job *Job) Exec() error {

	ctx := context.Background()

	workflow, err := job.toWorkflowExecution(ctx)
	if err != nil {
		return err
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

func (job *Job) toWorkflowExecution(ctx context.Context) (*daap.Execution, error) {

	env := []string{fmt.Sprintf("%s=%s", "HOTSUB_ROOT", HOTSUB_CONTAINERROOT)}
	for _, e := range job.Container.Envs {
		env = append(env, e.Pair())
	}

	workflow := &daap.Execution{Inspect: true, Env: env}

	if job.Type != CommonWorkflowLanguageJob {
		workflow.Script = job.Container.Script.Path
		return workflow, nil
	}

	// CWL
	// TODO: Uploading all "Includes" should be moved into "Fetch" step.
	for _, include := range job.Parameters.Includes {
		if err := job.upload(ctx, include); err != nil {
			return nil, err
		}
		envpair := fmt.Sprintf("%s=%s", include.Name, include.DeployedPath)
		env = append(env, envpair)
	}
	workflow.Env = env
	workflow.Inline = "cwltool ${CWL_FILE} ${CWL_PARAM_FILE}"

	return workflow, nil
}

func (job *Job) upload(ctx context.Context, include *Include) error {
	f, err := os.Open(include.LocalPath)
	if err != nil {
		return err
	}
	if err := job.Container.Workflow.Upload(ctx, f, HOTSUB_CONTAINERROOT); err != nil {
		return err
	}
	include.DeployedPath = filepath.ToSlash(filepath.Join(HOTSUB_CONTAINERROOT, filepath.Base(include.LocalPath)))
	return nil
}
