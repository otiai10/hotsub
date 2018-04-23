package core

import (
	"context"

	"golang.org/x/sync/semaphore"
)

// Run ...
func (job *Job) Run(ctx context.Context, shared *SharedData, sem *semaphore.Weighted) error {

	done := make(chan error)
	defer close(done)
	go job.run(shared, sem, done)

	// Destroy computing instance for this job anyway,
	// EVEN WHEN any of other job failed.
	defer job.Destroy()

	for {
		select {
		case err := <-done:
			// This "err" could be nil or non-nil.
			// Once non-nil error is returned, errgroup.WithContext would emit the cancellation
			// to all the "job.Run" by "ctx.Done()" channel.
			return err
		case <-ctx.Done():
			// If any of parallel "job.Run" has failed, the cancellation is notified by this channle.
			// Stop, destroy and return this "job.Run" as well.
			return ctx.Err()
		}
	}

}

func (job *Job) run(shared *SharedData, sem *semaphore.Weighted, done chan<- error) {

	jobctx := context.Background()
	defer jobctx.Done()

	sem.Acquire(jobctx, 1)
	if err := job.Create(); err != nil {
		done <- err
		return
	}
	sem.Release(1)

	if err := job.Construct(shared); err != nil {
		done <- err
		return
	}

	if err := job.Commit(); err != nil {
		done <- err
		return
	}

	done <- nil
	return
}
