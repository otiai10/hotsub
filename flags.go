package main

import (
	"github.com/urfave/cli"
)

// All flags listed here.
var flags = []cli.Flag{

	// Command Control
	VerboseFlag,
	LogDirectoryFlag,
	ConcurrencyFlag,

	// Debug
	KeepFlag,

	// Operation Contents
	ProviderFlag,
	TasksFlag,
	ImageFlag,
	ScriptFlag,

	SharedFlag,

	// Machine Specs
	// MinCoresFlag,
	// MinRAMFlag,
	DiskSizeFlag,

	// Amazon Web Service
	// AwsVPCFlag,
	AwsRegionFlag,
	AwsEC2InstanceType,
	AwsIAMInstanceProfile,

	//// Google Cloud Platform
	GoogleProjectFlag,
	GoogleZoneFlag,
	// GoogleBootDiskSizeFlag,
	// GooglePreEmptibleFlag,
	// GoogleScopesFlag,
	// GoogleKeepAlive,
	// GoogleAcceleratorTypeFlag,
}

// Debug flags

// VerboseFlag ...
var VerboseFlag = cli.BoolFlag{
	Name:  "verbose,v",
	Usage: `Print verbose log for operation.`,
}

// LogDirectoryFlag ...
var LogDirectoryFlag = cli.StringFlag{
	Name:  "log-dir",
	Usage: `Path to log directory where stdout/stderr log files will be placed (default: "${cwd}/logs/${time}")`,
}

// ConcurrencyFlag ...
var ConcurrencyFlag = cli.Int64Flag{
	Name:  "concurrency,C",
	Usage: `Concurrency for creating machines (≠ job running) // TODO: more documentation`,
	Value: 8,
}

// DryRunFlag ...
// var DryRunFlag = cli.BoolFlag{
// 	Name:  "dry-run",
// 	Usage: `Print the pipeline(s) that would be run and then exit. (default: false)`,
// }

// ProviderFlag ...
var ProviderFlag = cli.StringFlag{
	Name:  "provider,p",
	Usage: `Job service provider. Valid values are "aws" and "local."`,
	Value: "aws",
}

// TasksFlag ...
var TasksFlag = cli.StringFlag{
	Name:  "tasks",
	Usage: `Path to CSV of task parameters, expected to specify --env, --input, --input-recursive and --output-recursive. (required)`,
}

// ImageFlag ...
var ImageFlag = cli.StringFlag{
	Name:  "image",
	Usage: `Image name from Docker Hub or other Docker image service.`,
	Value: "ubuntu:14.04",
}

// ScriptFlag ...
var ScriptFlag = cli.StringFlag{
	Name:  "script",
	Usage: `Local path to a script to run inside the job's Docker container. (required)`,
}

// SharedFlag ...
var SharedFlag = cli.StringSliceFlag{
	Name:  "shared,S",
	Usage: `Shared data URL on cloud storage bucket. (e.g. s3://~)`,
}

// KeepFlag ...
var KeepFlag = cli.BoolFlag{
	Name:  "keep",
	Usage: `Keep instances created for computing event after everything gets done`,
}

// MinCoresFlag ...
// var MinCoresFlag = cli.UintFlag{
// 	Name:  "min-cores",
// 	Usage: `Minimum CPU cores for each job.`,
// 	Value: 1,
// }

// MinRAMFlag ...
// var MinRAMFlag = cli.Float64Flag{
// 	Name:  "min-ram",
// 	Usage: `Minimum RAM per job in GB.`,
// 	Value: 4,
// }

// DiskSizeFlag ...
var DiskSizeFlag = cli.UintFlag{
	Name:  "disk-size",
	Usage: `Size of data disk to attach for each job in GB.`,
	Value: 64,
}

//////////////////////////////////
// Flags for Amazon Web Service //
//////////////////////////////////

// AwsVPCFlag ...
// var AwsVPCFlag = cli.StringFlag{
// 	Name:  "aws-vpc",
// 	Usage: `AWS VPC ID in which AmazonEC2 instances would be launched`,
// }

// AwsRegionFlag ...
var AwsRegionFlag = cli.StringFlag{
	Name:  "aws-region",
	Usage: `AWS region name in which AmazonEC2 instances would be launched`,
	Value: "ap-northeast-1",
}

// AwsEC2InstanceType ...
var AwsEC2InstanceType = cli.StringFlag{
	Name:  "aws-ec2-instance-type",
	Usage: `AWS EC2 instance type. If specified, all --min-cores and --min-ram would be ignored.`,
	Value: "t2.micro",
}

// AwsIAMInstanceProfile ...
var AwsIAMInstanceProfile = cli.StringFlag{
	Name:  "aws-iam-instance-profile",
	Usage: `AWS instance profile from your IAM roles.`,
}

/////////////////////////////////////
// Flags for Google Cloud Platform //
/////////////////////////////////////

// GoogleProjectFlag ...
var GoogleProjectFlag = cli.StringFlag{
	Name:  "google-project",
	Usage: "Project ID for GCP",
}

// GoogleZoneFlag ...
var GoogleZoneFlag = cli.StringFlag{
	Name:  "google-zone",
	Usage: "GCP service zone name",
	Value: "asia-northeast1-a",
}
