package main

import (
	"context"
	"fmt"
	"net/url"
	"path/filepath"
	"sync"

	"github.com/otiai10/daap"
)

// Finalize ...
func (h *Handler) Finalize(ctx context.Context, c *daap.Container, job *Job) error {
	wg := new(sync.WaitGroup)
	for _, dest := range job.Task.OutputRecursive {
		wg.Add(1)
		h.uploadOutDirToCloud(ctx, c, job, dest, wg)
	}
	wg.Wait()
	return job.Error
}

func (h *Handler) uploadOutDirToCloud(ctx context.Context, c *daap.Container, job *Job, dest string, wg *sync.WaitGroup) {
	defer wg.Done()
	u, err := url.Parse(dest)
	if err != nil {
		job.Errorf("failed to parse destination url: %v", err)
		return
	}
	outdir := filepath.Join("/tmp", u.Path)
	execution := &daap.Execution{
		Inline: "/lifecycle/upload.sh",
		Env: []string{
			fmt.Sprintf("%s=%s", "SOURCE", outdir),
			fmt.Sprintf("%s=%s", "DEST", dest),
		},
		Inspect: true,
	}
	stream, err := c.Exec(ctx, execution)
	if err != nil {
		job.Errorf("failed to execute finalize: %v", err)
		return
	}
	for payload := range stream {
		job.Logf("[FINALIZE] &%d> %s", payload.Type, string(payload.Data))
	}
	if execution.ExitCode != 0 {
		job.Errorf("failed to upload output file `%s` with status code %d, please check output with --verbose option", dest, execution.ExitCode)
		return
	}
	job.Logf("[FINALIZE] Successfully uploaded: %v", dest)
	return
}
