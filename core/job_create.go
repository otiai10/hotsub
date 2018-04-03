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
func (job *Job) Create(shared *SharedData) error {

	job.Logf("[CREATE]\tCreating computing instance for this job...")

	spec := *job.Machine.Spec
	job.Identity.Name = fmt.Sprintf("%s-%04d", job.Identity.Prefix, job.Identity.Index)
	spec.Name = job.Identity.Name
	instance, err := dkmachine.Create(&spec)
	job.Machine.Instance = instance
	if err != nil {
		return err
	}

	if len(shared.Inputs) != 0 {
		if err := shared.CreateNFSVolumeOn(job.Machine.Instance); err != nil {
			return err
		}
		job.addContainerEnv(shared.Envs()...)
	}

	eg := new(errgroup.Group)
	eg.Go(func() error { return job.wakeupRoutineContainer() })
	eg.Go(func() error { return job.wakeupWorkflowContainer(shared) })
	return eg.Wait()
}

// wakeupRoutineContainer wakes the routine container up.
func (job *Job) wakeupRoutineContainer() error {

	job.Logf("[CREATE]\tSetting up routine container inside the computing instance...")

	ctx := context.Background()
	img := "awsub/lifecycle"
	container := daap.NewContainer(img, job.Machine.Instance)

	progress, err := container.PullImage(ctx)
	if err != nil {
		return err
	}
	job.drain(progress)

	err = container.Create(ctx, daap.CreateConfig{
		Host: &dockercontainer.HostConfig{
			Mounts:     []mount.Mount{daap.Bind(AWSUB_HOSTROOT, AWSUB_CONTAINERROOT)},
			Privileged: true,
		},
	})
	if err != nil {
		return err
	}

	job.Container.Routine = container

	return container.Start(ctx)
}

// wakeupWorkflowContainer wakes the user-defined workflow container up.
func (job *Job) wakeupWorkflowContainer(shared *SharedData) error {

	job.Logf("[CREATE]\tSetting up workflow container inside the computing instance...")

	ctx := context.Background()
	container := daap.NewContainer(job.Container.Image.Name, job.Machine.Instance)

	progress, err := container.PullImage(ctx)
	if err != nil {
		return err
	}
	job.drain(progress)

	mounts := []mount.Mount{
		daap.Bind(AWSUB_HOSTROOT, AWSUB_CONTAINERROOT),
	}
	if shared.Volume != nil && shared.Volume.Name != "" {
		mounts = append(
			mounts,
			daap.VolumeByName(shared.Volume.Name, AWSUB_CONTAINERROOT+"/"+AWSUB_SHARED_DIR),
			// daap.MountVolume(shared.Volume.Name, AWSUB_CONTAINERROOT+"/"+AWSUB_SHARED_DIR),
		)
	}

	err = container.Create(ctx, daap.CreateConfig{
		Host: &dockercontainer.HostConfig{
			Mounts: mounts,
		},
	})
	if err != nil {
		return err
	}

	job.Container.Workflow = container

	return container.Start(ctx)
}

func (job *Job) drain(ch <-chan daap.ImagePullResponsePayload) {
	for range ch {
		// fmt.Printf(".")
	}
	// fmt.Printf("\n")
}
