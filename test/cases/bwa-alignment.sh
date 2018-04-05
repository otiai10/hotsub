#!/bin/bash

PROJROOT=$(dirname $(dirname $(cd $(dirname $0) && pwd)))

set -e -v
awsub \
    --tasks ${PROJROOT}/test/bwa-alignment/bwa-alignment.csv \
    --script ${PROJROOT}/test/bwa-alignment/main.sh \
    --image otiai10/bwa \
    --aws-ec2-instance-type t2.large \
    --aws-iam-instance-profile testtest \
    --shared REFERENCE=s3://awsub-test-otiai10/bwa-example/reference/GRCh37 \
    --verbose
