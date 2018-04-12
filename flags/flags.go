package flags

import (
	"github.com/urfave/cli"
)

var verbose = cli.BoolFlag{
	Name:  "verbose,v",
	Usage: `Print verbose log for operation.`,
}

var logDirectory = cli.StringFlag{
	Name:  "log-dir",
	Usage: `Path to log directory where stdout/stderr log files will be placed (default: "${cwd}/logs/${time}")`,
}

var concurrency = cli.Int64Flag{
	Name:  "concurrency,C",
	Usage: `Concurrency for creating machines (≠ job running) // TODO: more documentation`,
	Value: 8,
}

// DryRun ...
// var DryRun = cli.BoolFlag{
// 	Name:  "dry-run",
// 	Usage: `Print the pipeline(s) that would be run and then exit. (default: false)`,
// }

var provider = cli.StringFlag{
	Name:  "provider,p",
	Usage: `Job service provider. Valid values are "aws" and "local."`,
	Value: "aws",
}

// tasks ...
var tasks = cli.StringFlag{
	Name:  "tasks",
	Usage: `Path to CSV of task parameters, expected to specify --env, --input, --input-recursive and --output-recursive. (required)`,
}

// image ...
var image = cli.StringFlag{
	Name:  "image",
	Usage: `Image name from Docker Hub or other Docker image service.`,
	Value: "ubuntu:14.04",
}

// script ...
var script = cli.StringFlag{
	Name:  "script",
	Usage: `Local path to a script to run inside the job's Docker container. (required)`,
}

// shared ...
var shared = cli.StringSliceFlag{
	Name:  "shared,S",
	Usage: `Shared data URL on cloud storage bucket. (e.g. s3://~)`,
}

var keep = cli.BoolFlag{
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

// disksize ...
var disksize = cli.UintFlag{
	Name:  "disk-size",
	Usage: `Size of data disk to attach for each job in GB.`,
	Value: 64,
}
