package core

import (
	"fmt"

	"golang.org/x/sync/errgroup"
)

// Create ...
func (component *Component) Create() error {

	// YAGNI: multiple SDIs for computing nodes
	if len(component.SharedData.Inputs) != 0 {
		// TODO: Use component.Logger
		fmt.Printf("[ROOT][CREATE]\tCreating Shared Data Instance...")
		if err := component.SharedData.Create(); err != nil {
			return err
		}
	}

	eg := new(errgroup.Group)

	for i, job := range component.Jobs {

		j := job
		j.Identity.Prefix = component.Identity.Name
		j.Identity.Index = i
		j.Machine.Spec = component.Machine.Spec

		if component.JobLoggerer != nil {
			j.Report.Log = component.JobLoggerer.Logger(j)
		}

		eg.Go(func() error {
			return j.Create(component.SharedData)
		})

		// Delegate runtimes
		j.Container.Image = component.Runtime.Image
		j.Container.Script = component.Runtime.Script
	}

	return eg.Wait()
}
