package core

import (
	"context"

	"golang.org/x/sync/semaphore"
)

// Run ...
func (job *Job) Run(ctx context.Context, shared *SharedData, sem *semaphore.Weighted) error {

	sem.Acquire(ctx, 1)
	if err := job.Create(); err != nil {
		return err
	}
	sem.Release(1)

	defer job.Destroy()

	if err := job.Construct(shared); err != nil {
		return err
	}

	if err := job.Commit(); err != nil {
		return err
	}

	return nil
}
