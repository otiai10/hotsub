#!/bin/bash

set -e

PROJROOT=$(dirname $(dirname $(cd $(dirname $0) && pwd)))

awsub \
    --tasks ${PROJROOT}/test/data/wordcount.csv \
    --script ${PROJROOT}/examples/wordcount/main.sh \
    --aws-iam-instance-profile testtest \
    --verbose
