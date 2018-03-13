# awsub [![Build Status](https://travis-ci.org/otiai10/awsub.svg?branch=master)](https://travis-ci.org/otiai10/awsub)

The simple batch job driver on AWS.

```sh
awsub \
  --image friend1ws/star-alignment \
  --script ./test/scripts/star-alignment.sh \
  --tasks ./test/tasks/star-alignment-tasks.csv \
  --aws-ec2-instance-type t2.2xlarge
```

# Installation

There are 3 options to install `awsub` command.

## 1. Download binary directly

Check releases here https://github.com/otiai10/awsub/releases and choose the binary for your OS.

## 2. go get

If you have `go` and `${GOPATH}/bin` is binded to `PATH`

```
go get -u github.com/otiai10/awsub
```

## 3. Clone and build

needs `go` as well

```
git clone git@github.com:otiai10/awsub.git
cd awsub
go install .
```

# Quick Guide

Once you have `awsub` installed, just hit following command

```
awsub quickguide
```

and you can see what you need.

# Contact

To make it transparent, ask any question from this link.

https://github.com/otiai10/awsub/issues
