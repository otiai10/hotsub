#!/bin/bash

set -e

PROJROOT=$(dirname $(dirname $(cd $(dirname $0) && pwd)))

# This control script fails intentionally.
# This test case checks if "awsub" can detect script error inside the container.
awsub \
    --tasks ${PROJROOT}/test/wordcount/wordcount.gcp.csv \
    --script ${PROJROOT}/test/wordcount/main.sh \
    --provider google \
    --google-zone asia-northeast1-a \
    --google-project genomondevel1 \
    --verbose
