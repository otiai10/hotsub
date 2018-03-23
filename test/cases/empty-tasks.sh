#!/bin/bash

set -o errexit
set -o verbose

PROJROOT=$(dirname $(dirname $(cd $(dirname $0) && pwd)))

awsub \
    --tasks ${PROJROOT}/test/wordcount/wordcount.empty.csv \
    --script ${PROJROOT}/test/wordcount/main.sh \
    --aws-iam-instance-profile testtest \
    --verbose
