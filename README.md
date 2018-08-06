# hotsub [![Build Status](https://travis-ci.org/otiai10/hotsub.svg?branch=master)](https://travis-ci.org/otiai10/hotsub)

The simple batch job driver on AWS and GCP. (Azure, OpenStack are coming soon)

```sh
hotsub run \
  --image friend1ws/star-alignment \
  --script ./test/scripts/star-alignment.sh \
  --tasks ./test/tasks/star-alignment-tasks.csv \
  --aws-ec2-instance-type t2.2xlarge
```

<img src="https://user-images.githubusercontent.com/931554/42975469-414f4878-8bf7-11e8-80fb-ad35d311fb6c.png" width="50%" />

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
