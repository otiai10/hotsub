package core

import (
	"fmt"
	"time"

	"github.com/otiai10/dkmachine/v0/dkmachine"
)

// Create creates physical machine and wake the required containers up.
// In most cases, containers with awsub/lifecycle and user defined image are required.
func (job *Job) Create() error {

	job.Lifetime(CREATE, "Creating computing instance for this job...")

	return job.create(0)
}

// create can be called recursively.
func (job *Job) create(retry int) error {

	if retry >= CreateMaxRetry {
		return fmt.Errorf("max retry of creating machine exceeded: failed %d times", CreateMaxRetry)
	}

	spec := *job.Machine.Spec
	job.Identity.Name = fmt.Sprintf("%s-%04d", job.Identity.Prefix, job.Identity.Index)
	spec.Name = job.Identity.Name
	instance, err := dkmachine.Create(&spec)
	job.Machine.Instance = instance

	// Succeeded!
	if err == nil {
		return nil
	}

	// Clean up for retry
	if err := instance.Remove(); err != nil {
		return fmt.Errorf("failed to clean up machine for retry: %v", err)
	}

	job.Lifetime(CREATE, "Retrying creating an instance for this job...")
	time.Sleep(time.Duration(retry*40) * time.Second)
	return job.create(retry + 1)
}
