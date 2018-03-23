#!/bin/bash

set -o errexit
set -o verbose

PROJROOT=$(dirname $(dirname $(cd $(dirname $0) && pwd)))

# This control script fails intentionally.
# This test case checks if "awsub" can detect script error inside the container.
awsub \
    --tasks ${PROJROOT}/test/control/example.csv \
    --script ${PROJROOT}/test/control/main.sh \
    --verbose
