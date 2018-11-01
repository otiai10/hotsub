# hotsub [![Build Status](https://travis-ci.org/otiai10/hotsub.svg?branch=master)](https://travis-ci.org/otiai10/hotsub)

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

# Installation

Check **[Getting Started](https://hotsub.github.io/getting-started)** on **[GitHub Pages](https://hotsub.github.io)**

# Contact

To make it transparent, ask any question from this link.

https://github.com/otiai10/hotsub/issues
