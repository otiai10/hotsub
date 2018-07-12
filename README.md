# hotsub [![Build Status](https://travis-ci.org/otiai10/hotsub.svg?branch=master)](https://travis-ci.org/otiai10/hotsub)

The simple batch job driver on AWS.

```sh
hotsub \
  --image friend1ws/star-alignment \
  --script ./test/scripts/star-alignment.sh \
  --tasks ./test/tasks/star-alignment-tasks.csv \
  --aws-ec2-instance-type t2.2xlarge
```

# Installation

Check releases here https://github.com/otiai10/awsub/releases and choose the binary for your OS.

# Quick Guide

Once you have `hotsub` installed, just hit following command

```
hotsub quickguide
```

and you can see what you need.

# Contact

To make it transparent, ask any question from this link.

https://github.com/otiai10/hotsub/issues
