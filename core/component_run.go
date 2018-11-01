package core

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

// Run executes all the jobs recursively.
// The concurrency of creating machines is managed here.
func (component *Component) Run(ctx context.Context) error {

	if len(component.Jobs) == 0 {
		return nil
	}

	eg, groupctx := errgroup.WithContext(ctx)
	concurrency := semaphore.NewWeighted(component.Concurrency)

	for i, job := range component.Jobs {

		// Ensure identity of this job.
		j := job
		j.Identity.Prefix = component.Identity.Name
		j.Identity.Index = i
		j.Identity.Name = fmt.Sprintf("%s.%d", j.Identity.Prefix, i)

		// Delegate specification of this job.
		j.Machine.Spec = component.Machine.Spec
		j.Container.Image = component.Runtime.Image
		j.Container.Script = component.Runtime.Script

		// Attach logger for this job.
		if err := component.loggerForJob(j); err != nil {
			return err
		}

		// Merge common parameters to each job.
		// TODO: Refactor, such as job.Parameters.Merge(common)
		j.Parameters.Envs = append(j.Parameters.Envs, component.CommonParameters.Envs...)

		// FIXME: Throttle API request to avoid "too many requests" error
		const AWS_API_REQUEST_LIMIT = 60
		if component.Concurrency >= AWS_API_REQUEST_LIMIT {
			time.Sleep(500 * time.Millisecond)
		}

		// Execute main.
		eg.Go(func() error {
			return j.Run(groupctx, component.SharedData, concurrency)
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}
