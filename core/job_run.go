package core

import (
	"context"

	"golang.org/x/sync/semaphore"
)

// Run ...
func (job *Job) Run(ctx context.Context, shared *SharedData, sem *semaphore.Weighted) error {

	done := job.run(shared, sem)

	// Destroy the instance independently,
	// NO MATTER if the neighbors got succeeded or failed.
	// See https://github.com/otiai10/hotsub/issues/115#issuecomment-417197674
	defer job.Destroy()

	for {
		select {
		case err := <-done:
			// This "err" could be nil or non-nil.
			// Once non-nil error is returned, errgroup.WithContext would emit the cancellation
			// to all the "job.Run" by "ctx.Done()" channel.
			return err
		case <-ctx.Done():
			/*
				// If any of parallel "job.Run" has failed, the cancellation is notified by this channle.
				// Stop, destroy and return this "job.Run" as well.
				return ctx.Err()
			*/

			// Do nothing. Don't let this job fail. Keep running.
			// See https://github.com/otiai10/hotsub/issues/101 for more detail.
		}
	}

}

func (job *Job) run(shared *SharedData, sem *semaphore.Weighted) <-chan error {

	jobctx := context.Background()
	defer jobctx.Done()

	done := make(chan error, 1)

	go func() {

		defer close(done)

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
	}()

	return done
}
