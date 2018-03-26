package core

import (
	"context"
	"fmt"

	dockercontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/otiai10/daap"
	"github.com/otiai10/dkmachine/v0/dkmachine"
	"golang.org/x/sync/errgroup"
)

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
