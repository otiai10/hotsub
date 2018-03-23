package platform

import "github.com/otiai10/dkmachine/v0/dkmachine"

// DefineMachineSpec is a factory layer to connect cli.Context to CreateOptions.
func DefineMachineSpec(ctx Context) (*dkmachine.CreateOptions, error) {
	opt := &dkmachine.CreateOptions{
		// AWS
		AmazonEC2Region:             ctx.String("aws-region"),
		AmazonEC2RootSize:           ctx.Int("disk-size"),
		AmazonEC2InstanceType:       ctx.String("aws-ec2-instance-type"),
		AmazonEC2IAMInstanceProfile: ctx.String("aws-iam-instance-profile"),
		// GCP
		GoogleProject:  ctx.String("google-project"),
		GoogleZone:     ctx.String("google-zone"),
		GoogleDiskSize: ctx.Int("disk-size"),
	}
	switch ctx.String("provider") {
	case "aws":
		opt.Driver = "amazonec2"
	case "gcp":
		opt.Driver = "google"
	default:
		opt.Driver = "amazonec2"
	}
	return opt, nil
}
