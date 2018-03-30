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

// SharedData ...
type SharedData struct {
	Spec      *dkmachine.CreateOptions
	Instance  *dkmachine.Machine
	Inputs    Inputs
	Root      string
	Container struct {
		Routine   *daap.Container
		NFSServer *daap.Container
	}
}

// Create ...
func (sd *SharedData) Create() error {

	instance, err := dkmachine.Create(sd.Spec)
	sd.Instance = instance
	if err != nil {
		return err
	}

	eg := new(errgroup.Group)
	eg.Go(sd.startNFS)
	eg.Go(sd.fetchAll)

	return eg.Wait()
}

func (sd *SharedData) fetchAll() error {

	ctx := context.Background()
	container := daap.NewContainer("awsub/lifecycle", sd.Instance)

	progress, err := container.PullImage(ctx)
	if err != nil {
		return nil
	}
	for range progress {
		fmt.Printf(".") // DEBUG: delete
	}

	err = container.Create(ctx, daap.CreateConfig{
		Host: &dockercontainer.HostConfig{
			Mounts: []mount.Mount{daap.MountVolume(AWSUB_MOUNTPOINT, AWSUB_MOUNTPOINT)},
		},
	})
	if err != nil {
		return err
	}

	if err := container.Start(ctx); err != nil {
		return err
	}

	sd.Container.Routine = container

	eg := new(errgroup.Group)

	for _, input := range sd.Inputs {
		i := input
		eg.Go(func() error { return sd.fetch(i) })
	}

	return eg.Wait()
}

func (sd SharedData) fetch(input *Input) error {

	ctx := context.Background()

	if err := input.Localize(AWSUB_CONTAINERROOT); err != nil {
		return err
	}

	fetch := &daap.Execution{
		Inline:  "/lifecycle/download.sh",
		Env:     input.EnvForFetch(),
		Inspect: true,
	}

	stream, err := sd.Container.Routine.Exec(ctx, fetch)
	if err != nil {
		return err
	}
	for payload := range stream {
		fmt.Printf("[SaredData] &%d> %s\n", payload.Type, payload.Text())
	}

	if fetch.ExitCode != 0 {
		return fmt.Errorf("fetch in SharedDataInstance exit with %d: %s", fetch.ExitCode, input.URL)
	}

	return nil
}

func (sd SharedData) startNFS() error {

	ctx := context.Background()
	container := daap.NewContainer("otiai10/nfs-server", sd.Instance)

	progress, err := container.PullImage(ctx)
	if err != nil {
		return nil
	}
	for range progress {
		fmt.Printf("#") // DEBUG: delete
	}

	err = container.Create(ctx, daap.CreateConfig{
		Host: &dockercontainer.HostConfig{
			Mounts:     []mount.Mount{daap.MountVolume(AWSUB_MOUNTPOINT, AWSUB_MOUNTPOINT)},
			Privileged: true,
			// NetworkMode: "host",
		},
		Container: &dockercontainer.Config{
			Env: []string{fmt.Sprintf("%s=%s", "MOUNTPOINT", AWSUB_MOUNTPOINT)},
		},
	})
	if err != nil {
		return err
	}

	if err := container.Start(ctx); err != nil {
		return err
	}

	return nil
}
