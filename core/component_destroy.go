package core

import (
	"time"

	"github.com/otiai10/stackerr"
)

// Destroy ...
func (component *Component) Destroy() error {

	errs := stackerr.New()

	// To make it sure to remove all the machines,
	// destroying computing instances should NOT be in part of job.Run.
	// This "component.Destroy" must be called even if any of "job.Run" get errored.
	for _, job := range component.Jobs {
		time.Sleep(500 * time.Millisecond)
		if err := job.Destroy(); err != nil {
			errs.Push(err)
		}
	}

	if component.SharedData.Instance != nil {
		if err := component.SharedData.Instance.Remove(); err != nil {
			errs.Push(err)
		}
	}

	return errs.Err()
}
