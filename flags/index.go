// Package flags defines all the flags awsub can accept.
// Any flags must be defined under this package.
package flags

import "github.com/urfave/cli"

// Index lists and exports all the flags so that awsub can use them.
var Index = []cli.Flag{
	verbose,
	logDirectory,
	concurrency,
	provider,
	tasks,
	image,
	script,
	shared,
	keep,
	disksize,
	env,

	// AWS
	awsRegion,
	awsEC2InstanceType,
	awsIAMInstanceProfile,

	// GCP
	googleProject,
	googleZone,
}
