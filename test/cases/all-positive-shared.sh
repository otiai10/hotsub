#!/bin/bash

PROJROOT=$(dirname $(dirname $(cd $(dirname $0) && pwd)))

set -e -v
awsub \
    --tasks ${PROJROOT}/test/wordcount/wordcount.csv \
    --script ${PROJROOT}/test/wordcount/main-shared.sh \
    --aws-iam-instance-profile testtest \
    --shared NYANCAT=s3://awsub-test-otiai10/speech/nyancat \
    --verbose
