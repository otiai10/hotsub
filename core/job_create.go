package core

import (
	"fmt"

	"github.com/otiai10/dkmachine/v0/dkmachine"
)

// Create creates physical machine and wake the required containers up.
// In most cases, containers with awsub/lifecycle and user defined image are required.
func (job *Job) Create() error {

	job.Lifetime(CREATE, "Creating computing instance for this job...")

	spec := *job.Machine.Spec
	job.Identity.Name = fmt.Sprintf("%s-%04d", job.Identity.Prefix, job.Identity.Index)
	spec.Name = job.Identity.Name
	instance, err := dkmachine.Create(&spec)
	job.Machine.Instance = instance
	if err != nil {
		return err
	}

	return nil
}
