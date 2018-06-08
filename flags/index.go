// Package flags defines all the flags awsub can accept.
// Any flags must be defined under this package.
package flags

import "github.com/urfave/cli"

// Index lists and exports all the flags so that awsub can use them.
var Index = []cli.Flag{
	Verbose,
	LogDirectory,
	Concurrency,
	Provider,
	Tasks,
	Image,
	Script,
	Shared,
	Keep,
	Env,
	Disksize,
	SharedDataDisksize,

	// AWS
	AwsRegion,
	AwsEC2InstanceType,
	AwsSharedInstanceType,

	// GCP
	GoogleProject,
	GoogleZone,
}
