#!/bin/bash

PROJROOT=$(dirname $(dirname $(cd $(dirname $0) && pwd)))

set -e -v

awsub \
    --tasks ${PROJROOT}/test/wordcount/wordcount.csv \
    --script ${PROJROOT}/test/wordcount/main.sh \
    --aws-iam-instance-profile testtest \
    --shared FOO=s3://hogefuga/piyo \
    --shared BAR=gs://foobar/baz \
    --verbose
