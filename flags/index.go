// Package flags defines all the flags hotsub can accept.
// Any flags must be defined under this package.
package flags

import "github.com/urfave/cli"

// Index lists and exports all the flags so that hotsub can use them.
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
	AwsVpcID,
	AwsSubnetID,

	// GCP
	GoogleProject,
	GoogleZone,

	// CommonWorkflowLanguage
	CWLFileFlag,
	CWLParamFlag,

	// WorkflowDescriptionLanguage
	WDLFileFlag,
	WDLJobFileFlag,

	// Not recommended
	IncludeFlag,
}
