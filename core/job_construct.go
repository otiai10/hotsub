package core

import (
	"context"

	dockercontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/otiai10/daap"
	"golang.org/x/sync/errgroup"
)

// Construct creates containers inside job instance.
func (job *Job) Construct(shared *SharedData) error {

	job.Lifetime("construct", "Constructing containers for this job...")

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

	job.Lifetime("construct", "Constructing routine container inside the computing instance...")

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

	job.Lifetime("construct", "Constructing workflow container inside the computing instance...")

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
