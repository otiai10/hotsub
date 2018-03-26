package core

import (
	"context"
	"fmt"
	"io"
	"time"

	dockercontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"

	"github.com/otiai10/daap"
	"github.com/otiai10/dkmachine/v0/dkmachine"
	"golang.org/x/sync/errgroup"
)

// NewJob ...
func NewJob(index int, prefix string) *Job {
	return &Job{
		Identity: Identity{
			Timestamp: time.Now().UnixNano(),
			Index:     index,
			Prefix:    prefix,
		},
		Container: JobContainer{
			Image:  &Image{},
			Script: &Script{},
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
	Container JobContainer

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

// Create creates physical machine and wake the required containers up.
// In most cases, containers with awsub/lifecycle and user defined image are required.
func (job *Job) Create() error {
	spec := *job.Machine.Spec
	job.Identity.Name = fmt.Sprintf("%s-%04d", job.Identity.Prefix, job.Identity.Index)
	spec.Name = job.Identity.Name
	instance, err := dkmachine.Create(&spec)
	if err != nil {
		return err
	}
	job.Machine.Instance = instance

	eg := new(errgroup.Group)
	eg.Go(job.wakeupRoutineContainer)
	eg.Go(job.wakeupWorkflowContainer)

	return eg.Wait()
}

// Destroy ...
func (job *Job) Destroy() error {
	if job.Machine.Instance == nil {
		return nil
	}
	return job.Machine.Instance.Remove()
}

// wakeupRoutineContainer wakes the routine container up.
func (job *Job) wakeupRoutineContainer() error {
	container, err := job.wakeupContainer("awsub/lifecycle")
	if err != nil {
		return err
	}
	job.Container.Routine = container
	return nil
}

// wakeupWorkflowContainer wakes the user-defined workflow container up.
func (job *Job) wakeupWorkflowContainer() error {
	container, err := job.wakeupContainer(job.Container.Image.Name)
	if err != nil {
		return err
	}
	job.Container.Workflow = container
	return nil
}

func (job *Job) wakeupContainer(img string) (*daap.Container, error) {

	ctx := context.Background()
	container := daap.NewContainer(img, job.Machine.Instance)

	progress, err := container.PullImage(ctx)
	if err != nil {
		return nil, err
	}
	job.drain(progress)

	err = container.Create(ctx, daap.CreateConfig{
		Host: &dockercontainer.HostConfig{
			Mounts: []mount.Mount{daap.Volume(AWSUB_HOSTROOT, AWSUB_CONTAINERROOT)},
		},
	})
	if err != nil {
		return nil, err
	}

	err = container.Start(ctx)
	return container, err
}

func (job *Job) drain(ch <-chan daap.ImagePullResponsePayload) {
	for range ch {
		// fmt.Printf(".")
	}
	// fmt.Printf("\n")
}

// Commit represents a main process of this job.
// The main process of this job consists of Fetch, Exec, and Push.
func (job *Job) Commit() error {

	if err := job.Fetch(); err != nil {
		return err
	}

	if err := job.Exec(); err != nil {
		return err
	}

	if err := job.Push(); err != nil {
		return err
	}

	return nil
}

// Push uploads result files to cloud storage services according to
// specified "output" URLs of tasks file.
func (job *Job) Push() error {
	return nil
}
