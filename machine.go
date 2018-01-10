package main

import (
	"fmt"

	"github.com/otiai10/dkmachine/v0/dkmachine"
)

func (h *Handler) generateMachineOption(task *Task) (*dkmachine.CreateOptions, error) {

	name := fmt.Sprintf("%s%02d", task.Prefix, task.Index)
	opt := &dkmachine.CreateOptions{
		Name: name,
	}
	opt.Driver = "amazonec2"
	opt.AmazonEC2Region = "ap-southeast-2"
	opt.AmazonEC2IAMInstanceProfile = "testtest"
	opt.AmazonEC2SecurityGroup = name

	return opt, nil
}
