package main

import (
	"fmt"

	"github.com/otiai10/dkmachine/v0/dkmachine"
)

func (h *Handler) generateMachineOption(task *Task) (*dkmachine.CreateOptions, error) {

	name := fmt.Sprintf("%s%02d", task.Prefix, task.Index)
	opt := &dkmachine.CreateOptions{Name: name}
	var err error
	switch h.ctx.String("provider") {
	case "aws":
		err = h.setupAWSMachineOption(opt)
	case "google":
		err = h.setupGCPMachineOption(opt)
	default:
		err = h.setupAWSMachineOption(opt)
	}
	if err != nil {
		return nil, err
	}

	return opt, nil
}

func (h *Handler) setupAWSMachineOption(opt *dkmachine.CreateOptions) error {

	opt.Driver = "amazonec2"

	opt.AmazonEC2RootSize = h.ctx.Int("disk-size")
	// e.g. "ap-southeast-2"
	opt.AmazonEC2Region = h.ctx.String("aws-region")
	// e.g. "t2.2xlarge"
	opt.AmazonEC2InstanceType = h.ctx.String("aws-ec2-instance-type")
	// e.g. "my-role"
	opt.AmazonEC2IAMInstanceProfile = h.ctx.String("aws-iam-instance-profile")

	opt.AmazonEC2SecurityGroup = opt.Name

	// FIXME: hard coding
	opt.AmazonEC2RequestSpotInstance = false

	return nil
}

func (h *Handler) setupGCPMachineOption(opt *dkmachine.CreateOptions) error {

	opt.Driver = "google"

	opt.GoogleProject = h.ctx.String("google-project")

	return nil
}
