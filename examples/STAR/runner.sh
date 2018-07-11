#!/bin/bash

CWD=$(cd $(dirname $0) && pwd)

set -v
# Hi, this is a note for `hotsub` users.
# This script is an example to run STAR RNA-seq alignment.
# The following command definitely represents how `hotsub` works.
# You might need to customize 2 things
#     1. The location of sample file described in CSV file.
#     2. The location of reference files specified by `--shared` flag.
#     3. IAM instance profile that you have.
# and you can have this example working for your own samples.
# Good luck!
hotsub \
  --tasks ${CWD}/star-alignment.csv \
  --shared REFERENCE=s3://hotsub/resources/reference/GRCh37.STAR-2.5.2a \
  --script ${CWD}/main.sh \
  --image friend1ws/star-alignment \
  --aws-ec2-instance-type m4.2xlarge \
  --disk-size 128 \
  --verbose
