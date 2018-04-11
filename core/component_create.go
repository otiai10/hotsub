package core

import (
	"context"
	"fmt"

	"golang.org/x/sync/semaphore"

	"golang.org/x/sync/errgroup"
)

// Create creates the instances for each Job.
func (component *Component) Create() error {

	eg := new(errgroup.Group)

	// YAGNI: multiple SDIs for computing nodes
	if len(component.SharedData.Inputs) != 0 {
		// TODO: Use component.Logger
		eg.Go(component.SharedData.Create)
	}

	ctx := context.Background()
	sem := semaphore.NewWeighted(4)

	for i, job := range component.Jobs {

		j := job
		j.Identity.Prefix = component.Identity.Name
		j.Identity.Index = i
		j.Identity.Name = fmt.Sprintf("%s.%d", j.Identity.Prefix, i)
		j.Machine.Spec = component.Machine.Spec

		if component.JobLoggerFactory != nil {
			logger, err := component.JobLoggerFactory.Logger(j)
			if err != nil {
				return err
			}
			j.Report.Log = logger
		}

		if err := sem.Acquire(ctx, 1); err != nil {
			return err
		}
		eg.Go(func() error {
			defer sem.Release(1)
			return j.Create()
		})

		// Delegate runtimes
		j.Container.Image = component.Runtime.Image
		j.Container.Script = component.Runtime.Script
	}

	return eg.Wait()
}
