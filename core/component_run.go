package core

import (
	"fmt"

	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

// Run executes all the jobs recursively.
// The concurrency of creating machines is managed here.
func (component *Component) Run() error {

	if len(component.Jobs) == 0 {
		return nil
	}

	eg := new(errgroup.Group)
	sem := semaphore.NewWeighted(component.Concurrency)

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

		// Execute main.
		eg.Go(func() error {
			return j.Run(component.SharedData, sem)
		})
	}

	return eg.Wait()
}

// Prepare ...
func (component *Component) Prepare() error {
	if len(component.SharedData.Inputs) != 0 {
		if err := component.SharedData.Create(); err != nil {
			return err
		}
	}
	return nil
}

func (component *Component) loggerForJob(job *Job) error {
	if component.JobLoggerFactory == nil {
		return nil
	}
	logger, err := component.JobLoggerFactory.Logger(job)
	if err != nil {
		return err
	}
	job.Report.Log = logger
	return nil
}
