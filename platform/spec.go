package platform

import (
	"strings"

	"github.com/otiai10/awsub/core"
	"github.com/otiai10/dkmachine/v0/dkmachine"
)

// DefineMachineSpec is a factory layer to connect cli.Context to CreateOptions.
func DefineMachineSpec(ctx Context) (*dkmachine.CreateOptions, error) {
	opt := &dkmachine.CreateOptions{
		// AWS
		AmazonEC2Region:             ctx.String("aws-region"),
		AmazonEC2RootSize:           ctx.Int("disk-size"),
		AmazonEC2InstanceType:       ctx.String("aws-ec2-instance-type"),
		AmazonEC2IAMInstanceProfile: ctx.String("aws-iam-instance-profile"),
		AmazonEC2SecurityGroup:      DefaultAWSSecurityGroupName,
		// GCP
		GoogleProject:  ctx.String("google-project"),
		GoogleZone:     ctx.String("google-zone"),
		GoogleDiskSize: ctx.Int("disk-size"),
		GoogleScopes: strings.Join([]string{
			"https://www.googleapis.com/auth/devstorage.read_write",
			"https://www.googleapis.com/auth/logging.write,https://www.googleapis.com/auth/monitoring.write",
		}, ","),
	}
	switch Provider(ctx.String("provider")) {
	case AWS:
		opt.Driver = string(AmazonEC2)
	case GCP:
		opt.Driver = string(Google)
	default:
		opt.Driver = string(AmazonEC2)
	}
	return opt, nil
}

// DefineSharedDataInstanceSpec ...
// It defines ALL the specifications for docker-machine,
// though it's a bit verbose ;)
func DefineSharedDataInstanceSpec(shared core.Inputs, ctx Context) (*dkmachine.CreateOptions, error) {
	opt := &dkmachine.CreateOptions{
		AmazonEC2Region:             ctx.String("aws-region"),
		AmazonEC2IAMInstanceProfile: ctx.String("aws-iam-instance-profile"),
		AmazonEC2SecurityGroup:      DefaultAWSSecurityGroupName,
		// {{{ TODO: Fix hard coding
		AmazonEC2InstanceType: "m4.2xlarge", // TODO: Fix hard coding
		AmazonEC2RootSize:     64,           // TODO: Fix hard coding
		// }}}
		Name: "Shared-Data-Instance",
	}
	switch Provider(ctx.String("provider")) {
	case AWS:
		opt.Driver = string(AmazonEC2)
	case GCP:
		opt.Driver = string(Google)
	default:
		opt.Driver = string(AmazonEC2)
	}
	return opt, nil
}
