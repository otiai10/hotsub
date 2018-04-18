package core

import (
	"time"

	"github.com/otiai10/daap"
	"github.com/otiai10/dkmachine/v0/dkmachine"
)

// NewJob ...
func NewJob(index int, prefix string) *Job {
	return &Job{
		Identity: Identity{
			Timestamp: time.Now().UnixNano(),
			Index:     index,
			Prefix:    prefix,
		},
		Parameters: &Parameters{},
		Container: &JobContainer{
			Image:  &Image{},
			Script: &Script{},
		},
		Report: &Report{},
	}
}

// Job represents a input/output/env set specified as an independent row of tasks file.
type Job struct {

	// Identity specifies the identity of this job.
	Identity Identity

	// Parameters specifies the parameters assigned to this job.
	// It is exactly what the corresponding row in tasks file is parsed to.
	Parameters *Parameters

	// Container spedifies the settings which is used the real execution runtime.
	Container *JobContainer

	Machine struct {
		Spec     *dkmachine.CreateOptions
		Instance *dkmachine.Machine
	}

	// Report ...
	Report *Report
}

// JobContainer ...
type JobContainer struct {
	// Envs shold have evetything translated from Parameters.
	Envs   []Env
	Image  *Image
	Script *Script

	// container ...
	Routine  *daap.Container
	Workflow *daap.Container
}
