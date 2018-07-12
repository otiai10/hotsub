#!/bin/bash

PROJROOT=$(dirname $(dirname $(cd $(dirname $0) && pwd)))

set -e -v
hotsub \
    --tasks ${PROJROOT}/test/bwa-alignment/bwa-alignment.csv \
    --script ${PROJROOT}/test/bwa-alignment/main.sh \
    --image otiai10/bwa \
    --aws-ec2-instance-type t2.large \
    --shared REFERENCE=s3://hotsub/resources/reference/GRCh37 \
    --env REFERENCE_FILE=GRCh37.fa \
    --verbose
