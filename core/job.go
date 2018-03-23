package core

import (
	"fmt"
	"io"
	"time"

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
	}
}

// Job represents a input/output/env set specified as an independent row of tasks file.
type Job struct {

	// Identity specifies the identity of this job.
	Identity Identity

	// Parameters specifies the parameters assigned to this job.
	// It is exactly what the corresponding row in tasks file is parsed to.
	Parameters struct {
		Inputs  Inputs
		Outputs Outputs
		Envs    []Env
	} `json:"parameters"`

	// Container spedifies the settings which is used the real execution runtime.
	Container struct {
		// Envs shold have evetything translated from Parameters.
		Envs   []Env
		Image  *Image
		Script *Script
		// Instance keeps the informations of physical machine, might be created by docker-machine.
	}

	Machine struct {
		Spec     *dkmachine.CreateOptions
		Instance *dkmachine.Machine
	}

	// Report ...
	Report *Report
}

// Report ...
type Report struct {
	Log struct {
		Writer io.Writer
	}
	Metrics struct {
		Writer io.Writer
	}
}

// Create ...
func (job *Job) Create() error {
	spec := *job.Machine.Spec
	spec.Name = fmt.Sprintf("%s-%04d", job.Identity.Prefix, job.Identity.Index)
	fmt.Printf("%+v\n", spec)
	instance, err := dkmachine.Create(&spec)
	job.Machine.Instance = instance
	return err
}

// Destroy ...
func (job *Job) Destroy() error {
	if job.Machine.Instance == nil {
		return nil
	}
	return job.Machine.Instance.Remove()
}
