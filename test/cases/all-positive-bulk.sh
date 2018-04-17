#!/bin/bash

PROJROOT=$(dirname $(dirname $(cd $(dirname $0) && pwd)))

set -e -v
awsub \
    --tasks ${PROJROOT}/test/wordcount/wordcount.bulk.csv \
    --script ${PROJROOT}/test/wordcount/main.sh \
    --aws-iam-instance-profile awsubtest \
    --log-dir /tmp \
    --verbose
