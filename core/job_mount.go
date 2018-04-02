package core

import (
	"context"
	"fmt"

	"github.com/otiai10/daap"
	"github.com/otiai10/dkmachine/v0/dkmachine"
)

// Mount ...
func (job *Job) Mount(shared *dkmachine.Machine) error {
	inline := fmt.Sprintf(
		"mkdir -p %[1]s && mount -t nfs %[2]s:/ %[1]s",
		AWSUB_CONTAINERROOT+"/"+AWSUB_SHARED_DIR,
		shared.Driver.PrivateIPAddress,
	)
	mnt := &daap.Execution{Inline: inline, Inspect: true}
	ctx := context.Background()
	stream, err := job.Container.Routine.Exec(ctx, mnt)
	if err != nil {
		return err
	}
	for payload := range stream {
		fmt.Printf("[MOUNT] &%d> %s\n", payload.Type, payload.Text())
	}

	if mnt.ExitCode != 0 {
		return fmt.Errorf("mount onto SharedDataInstance exit with %d: %s", mnt.ExitCode, inline)
	}
	return nil
}
