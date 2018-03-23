#!/bin/bash

set -o errexit
set -o verbose

PROJROOT=$(dirname $(dirname $(cd $(dirname $0) && pwd)))

awsub \
    --tasks ${PROJROOT}/test/wordcount/wordcount.csv \
    --script ${PROJROOT}/test/wordcount/main.sh \
    --verbose
    # --aws-iam-instance-profile testtest
