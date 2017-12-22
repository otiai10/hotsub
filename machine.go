package main

import (
	"fmt"

	"github.com/otiai10/dkmachine/v0/dkmachine"
)

func (h *Handler) generateMachineOption(task *Task) (*dkmachine.CreateOptions, error) {

	opt := &dkmachine.CreateOptions{
		Name: fmt.Sprintf("%s%02d", task.Prefix, task.Index),
	}
	opt.Driver = "amazonec2"
	opt.AmazonEC2Region = "ap-southeast-2"
	opt.AmazonEC2IAMInstanceProfile = "testtest"

	return opt, nil
}
