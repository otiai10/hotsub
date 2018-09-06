package core

import (
	"fmt"
	"time"

	"github.com/docker/machine/libmachine/check"
	"github.com/otiai10/dkmachine"
)

// Create creates physical machine and wake the required containers up.
// In most cases, containers with hotsub/routine and user defined image are required.
func (job *Job) Create() error {

	job.Lifetime(CREATE, "Creating computing instance for this job...")

	return job.create(0, false, nil)
}

// create can be called recursively.
func (job *Job) create(retry int, regenerateCerts bool, lasterror error) error {

	if retry >= CreateMaxRetry {
		return fmt.Errorf("max retry of creating machine exceeded: failed %d times with last error: %v", CreateMaxRetry, lasterror)
	}

	var err error
	if regenerateCerts {
		err = job.Machine.Instance.RegenerateCerts()
	} else {
		spec := *job.Machine.Spec
		job.Identity.Name = fmt.Sprintf("%s-%04d", job.Identity.Prefix, job.Identity.Index)
		spec.Name = job.Identity.Name
		job.Machine.Instance, err = dkmachine.Create(&spec)
	}

	// Succeeded!
	if err == nil {
		return nil
	}

	if _, ok := err.(check.ErrCertInvalid); ok {
		job.Lifetime(CREATE, "Regenerating certificates for this job after %d seconds. REASON: %T", (retry * 5), err)
		time.Sleep(time.Duration(retry*5) * time.Second)
		return job.create(retry+1, true, err)
	}

	// Clean up for retry
	if errOnRemove := job.Machine.Instance.Remove(); errOnRemove != nil {
		return fmt.Errorf("last error on create: %v: failed to clean up machine for retry: %v", err, errOnRemove)
	}

	job.Lifetime(CREATE, "Retrying instance creation for this job after %d seconds. REASON: %T", (retry * 5), err)
	time.Sleep(time.Duration(retry*5) * time.Second)
	return job.create(retry+1, false, err)
}
