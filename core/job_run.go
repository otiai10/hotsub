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

	for {
		select {
		case err := <-done:
			return err
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

	defer job.Destroy()

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
