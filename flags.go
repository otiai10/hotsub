package main

import "github.com/urfave/cli"

// All flags listed here.
var flags = []cli.Flag{
	ProviderFlag,
	NameFlag,
	TasksFlag,
	ImageFlag,
	DryRunFlag,
	CommandFlag,
	ScriptFlag,
	EnvFlag,
	LabelFlag,
	InputFlag,
	InputRecursiveFlag,
	OutputFlag,
	OutputRecursiveFlag,
	UserFlag,
	WaitFlag,
	PoleIntervalFlag,
	AfterFlag,
	SkipFlag,
	MinCoresFlag,
	MinRAMFlag,
	DiskSizeFlag,
	LoggingFlag,

	// Amazon Web Service
	AwsVPCFlag,
	AwsRegionFlag,

	// Google Cloud Platform
	GoogleProjectFlag,
	GoogleBootDiskSizeFlag,
	GooglePreEmptibleFlag,
	GoogleZonesFlag,
	GoogleScopesFlag,
	GoogleKeepAlive,
	GoogleAcceleratorTypeFlag,
}

// ProviderFlag ...
var ProviderFlag = cli.StringFlag{
	Name:  "provider",
	Usage: `Job service provider. Valid values are "aws" and "local" (local Docker execution). (default: aws)`,
	Value: "aws",
}

// NameFlag ...
var NameFlag = cli.StringFlag{
	Name:  "name",
	Usage: `Name for pipeline. Defaults to the script name or first token of the --command if specified. (default: "")`,
}

// TasksFlag ...
var TasksFlag = cli.StringFlag{
	Name:  "tasks",
	Usage: `Path to TSV of task parameters. Each column can specify an --env, --input, or --output variable, and each line specifies the values of those variables for a separate task. Optionally specify tasks in the file to submit. Can take the form "m", "m-", or "m-n" where m and n are task numbers. (default: None)`,
}

// ImageFlag ...
var ImageFlag = cli.StringFlag{
	Name:  "image",
	Usage: `Image name from Docker Hub, Google Container Repository, or other Docker image service. The pipeline must have READ access to the image. (default: ubuntu:14.04)`,
	Value: "ubuntu:14.04",
}

// DryRunFlag ...
var DryRunFlag = cli.BoolFlag{
	Name:  "dry-run",
	Usage: `Print the pipeline(s) that would be run and then exit. (default: False)`,
}

// CommandFlag ...
var CommandFlag = cli.StringFlag{
	Name:  "command",
	Usage: `Command to run inside the job's Docker container (default: None)`,
}

// ScriptFlag ...
var ScriptFlag = cli.StringFlag{
	Name:  "script",
	Usage: `Local path to a script to run inside the job's Docker container. (default: None)`,
}

// EnvFlag ...
var EnvFlag = cli.StringSliceFlag{
	Name:  "env",
	Usage: `Environment variables for the script's execution environment (default: [])`,
}

// LabelFlag ...
var LabelFlag = cli.StringSliceFlag{
	Name:  "label",
	Usage: `Labels to associate to the job. (default: [])`,
}

// InputFlag ...
var InputFlag = cli.StringSliceFlag{
	Name:  "input",
	Usage: `Input path arguments to localize into the script's execution environment (default: [])`,
}

// InputRecursiveFlag ...
var InputRecursiveFlag = cli.StringSliceFlag{
	Name:  "input-recursive",
	Usage: `Input path arguments to localize recursively into the script's execution environment (default: [])`,
}

// OutputFlag ...
var OutputFlag = cli.StringSliceFlag{
	Name:  "output",
	Usage: `Output path arguments to de-localize from the script's execution environment (default: [])`,
}

// OutputRecursiveFlag ...
var OutputRecursiveFlag = cli.StringSliceFlag{
	Name:  "output-recursive",
	Usage: `Output path arguments to de-localize recursively from the script's execution environment (default: [])`,
}

// UserFlag ...
var UserFlag = cli.StringFlag{
	Name:  "user",
	Usage: `User submitting the awsub job, defaults to the current OS user. (default: None)`,
}

// WaitFlag ...
var WaitFlag = cli.BoolFlag{
	Name:  "wait",
	Usage: `Wait for the job to finish all its tasks. (default: False)`,
}

// PoleIntervalFlag ...
var PoleIntervalFlag = cli.UintFlag{
	Name:  "pole-interval",
	Usage: `Polling interval (in seconds) for checking job status when --wait or --after are set. (default: 20)`,
	Value: 20,
}

// AfterFlag ...
var AfterFlag = cli.StringSliceFlag{
	Name:  "after",
	Usage: `Job ID(s) to wait for before starting this job. (default: [])`,
}

// SkipFlag ...
var SkipFlag = cli.BoolFlag{
	Name:  "skip",
	Usage: `Do not submit the job if all output specified using the --output and --output-recursive parameters already exist. Note that wildcard and recursive outputs cannot be strictly verified. See the documentation for details. (default: False)`,
}

// MinCoresFlag ...
var MinCoresFlag = cli.UintFlag{
	Name:  "min-cores",
	Usage: `Minimum CPU cores for each job (default: 1)`,
	Value: 1,
}

// MinRAMFlag ...
var MinRAMFlag = cli.Float64Flag{
	Name:  "min-ram",
	Usage: `Minimum RAM per job in GB (default: 4)`,
	Value: 4,
}

// DiskSizeFlag ...
var DiskSizeFlag = cli.UintFlag{
	Name:  "disk-size",
	Usage: `Size (in GB) of data disk to attach for each job (default: 200)`,
	Value: 200,
}

// LoggingFlag ...
var LoggingFlag = cli.StringFlag{
	Name:  "logging",
	Usage: `Cloud Storage path to send logging output (either a folder, or file ending in ".log") (default: None)`,
}

/**
 * Flags for Amazon Web Service
 */

// AwsVPCFlag ...
var AwsVPCFlag = cli.StringFlag{
	Name: "aws-vpc",
}

// AwsRegionFlag ...
var AwsRegionFlag = cli.StringFlag{
	Name: "aws-region",
}

/**
 * Flags for Google Cloud Platform
 */

// GoogleProjectFlag ...
var GoogleProjectFlag = cli.StringFlag{
	Name:  "google-project",
	Usage: `Cloud project ID in which to run the pipeline (default: None)`,
}

// GoogleBootDiskSizeFlag ...
var GoogleBootDiskSizeFlag = cli.UintFlag{
	Name:  "google-boot-disk-size",
	Usage: `Size (in GB) of the boot disk (default: 10)`,
	Value: 10,
}

// GooglePreEmptibleFlag ...
var GooglePreEmptibleFlag = cli.BoolFlag{
	Name:  "google-preemptible",
	Usage: `Use a preemptible VM for the job (default: False)`,
}

// GoogleZonesFlag ...
var GoogleZonesFlag = cli.StringSliceFlag{
	Name:  "google-zones",
	Usage: `List of Google Compute Engine zones. (default: None)`,
}

// GoogleScopesFlag ...
var GoogleScopesFlag = cli.StringSliceFlag{
	Name:  "google-scopes",
	Usage: `Space-separated scopes for GCE instances. (default:`,
	Value: &cli.StringSlice{"https://www.googleapis.com/auth/bigquery"},
}

// GoogleKeepAlive ...
var GoogleKeepAlive = cli.UintFlag{
	Name:  "google-keepalive",
	Usage: `Time (in seconds) to keep a tasks's virtual machine (VM) running after a localization, docker command, or delocalization failure. Allows for connecting to the VM for debugging. Default is 0; maximum allowed value is 86400 (1 day). (default: None)`,
}

// GoogleAcceleratorTypeFlag ...
var GoogleAcceleratorTypeFlag = cli.StringFlag{
	Name:  "google-accelerator-type",
	Usage: `The Compute Engine accelerator type. By specifying this parameter, you will download and install the following third-party software onto your job's Compute Engine instances: NVIDIA(R) Tesla(R) drivers and NVIDIA(R) CUDA toolkit. Please see https://cloud.google.com/compute/docs/gpus/ for supported GPU types and https://cloud.google.com/genom ics/reference/rest/v1alpha2/pipelines#pipelineresources for more details. (default: None)`,
}
