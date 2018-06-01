#!/bin/bash

set -e -v
awsub \
    --tasks ./bwa-alignment.csv \
    --script ./main.sh \
    --image otiai10/bwa \
    --aws-ec2-instance-type t2.large \
    --aws-iam-instance-profile awsubtest \
    --verbose
