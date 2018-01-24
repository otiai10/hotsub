# awsub [![Build Status](https://travis-ci.org/otiai10/awsub.svg?branch=master)](https://travis-ci.org/otiai10/awsub)

```sh
awsub \
  --image friend1ws/star-alignment \
  --script ./test/scripts/star-alignment.sh \
  --tasks ./test/tasks/star-alignment-tasks.csv \
  --aws-ec2-instance-type t2.2xlarge
```

# Installation

If you have `go` and `GOPATH/bin` avialble for `PATH`, just hit

```
% go get -u github.com/otiai10/awsub
% awsub quickguide
```
