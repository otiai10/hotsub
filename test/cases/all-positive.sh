#!/bin/bash

PROJROOT=$(dirname $(dirname $(cd $(dirname $0) && pwd)))

set -e -v
awsub \
    --tasks ${PROJROOT}/test/wordcount/wordcount.csv \
    --script ${PROJROOT}/test/wordcount/main.sh \
    --env FOO="This is foo, specified as a common env" \
    --aws-iam-instance-profile awsubtest \
    --log-dir /tmp \
    --verbose
