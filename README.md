# hotsub [![Build Status](https://travis-ci.org/otiai10/hotsub.svg?branch=master)](https://travis-ci.org/otiai10/hotsub) [![Paper Status](http://joss.theoj.org/papers/f1e4470e4831caa4252427cec8c009a8/status.svg)](http://joss.theoj.org/papers/f1e4470e4831caa4252427cec8c009a8)

The simple batch job driver on AWS and GCP. (Azure, OpenStack are coming soon)

```sh
hotsub run \
  --script ./star-alignment.sh \
  --tasks ./star-alignment-tasks.csv \
  --image friend1ws/star-alignment \
  --aws-ec2-instance-type t2.2xlarge \
  --verbose
```

It will

- execute workflow described in `star-alignment.sh`
- for each samples specified in `star-alignment.csv`
- in `friend1ws/star-alignment` docker containers
- on EC2 instances of type `t2.2xlarge`

and automatically upload the output files to S3 and clean up EC2 instances after all.

See **[Documentation](https://hotsub.github.io/)** for more details.

# Why you use `hotsub`

There are 3 points why `hotsub` is made and why you use it

1. **No-need to setup your cloud on web consoles:**
    - Since `hotsub` uese pure EC2 or GCE instances, you don't have to configure AWS Batch nor Dataflow on messy web consoles
2. **Multi-platforms with the same interface of command line:**
    - You can switch AWS and GCP as you like only with `--provider` option of `run` command (of course you need to have credentials on your local machine)
3. **ExTL framework available:**
    - In some cases of bio-informatics, the problem is how to handle common and huge refrence genome. `hotsub` suggests and implements <u>[`ExTL` framework](https://hotsub.github.io/etl-and-extl)</u>.

# Installation

Check **[Getting Started](https://hotsub.github.io/getting-started)** on **[GitHub Pages](https://hotsub.github.io)**

# Commands

```sh
NAME:
   hotsub - command line to run batch computing on AWS and GCP with the same interface

USAGE:
   hotsub [global options] command [command options] [arguments...]

VERSION:
   0.10.0

DESCRIPTION:
   Open-source command-line tool to run batch computing tasks and workflows on backend services such as Amazon Web Services.

COMMANDS:
     run       Run your jobs on cloud with specified input files and any parameters
     init      Initialize CLI environment on which hotsub runs
     template  Create a template project of hotsub
     help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -V  print the version
```

## Available options for `run` command


```sh
% hotsub run -h
NAME:
   hotsub run - Run your jobs on cloud with specified input files and any parameters

USAGE:
   hotsub run [command options] [arguments...]

DESCRIPTION:
   Run your jobs on cloud with specified input files and any parameters

OPTIONS:
   --verbose, -v                     Print verbose log for operation.
   --log-dir value                   Path to log directory where stdout/stderr log files will be placed (default: "${cwd}/logs/${time}")
   --concurrency value, -C value     Throttle concurrency number for running jobs (default: 8)
   --provider value, -p value        Job service provider, either of [aws, gcp, vbox, hyperv] (default: "aws")
   --tasks value                     Path to CSV of task parameters, expected to specify --env, --input, --input-recursive and --output-recursive. (required)
   --image value                     Image name from Docker Hub or other Docker image service. (default: "ubuntu:14.04")
   --script value                    Local path to a script to run inside the workflow Docker container. (required)
   --shared value, -S value          Shared data URL on cloud storage bucket. (e.g. s3://~)
   --keep                            Keep instances created for computing event after everything gets done
   --env value, -E value             Environment variables to pass to all the workflow containers
   --disk-size value                 Size of data disk to attach for each job in GB. (default: 64)
   --shareddata-disksize value       Disk size of shared data instance (in GB) (default: 64)
   --aws-region value                AWS region name in which AmazonEC2 instances would be launched (default: "ap-northeast-1")
   --aws-ec2-instance-type value     AWS EC2 instance type. If specified, all --min-cores and --min-ram would be ignored. (default: "t2.micro")
   --aws-shared-instance-type value  Shared Instance Type on AWS (default: "m4.4xlarge")
   --aws-vpc-id value                VPC ID on which computing VMs are launched
   --aws-subnet-id value             Subnet ID in which computing VMs are launched
   --google-project value            Project ID for GCP
   --google-zone value               GCP service zone name (default: "asia-northeast1-a")
   --cwl value                       CWL file to run your workflow
   --cwl-job value                   Parameter files for CWL
   --wdl value                       WDL file to run your workflow
   --wdl-job value                   Parameter files for WDL
   --include value                   Local files to be included onto workflow container
```

# Contact

To make it transparent, ask any question from this link.

https://github.com/otiai10/hotsub/issues
