#!/bin/bash

PROJROOT=$(dirname $(dirname $(cd $(dirname $0) && pwd)))

set -e -v
hotsub run \
    --tasks ${PROJROOT}/test/wordcount/wordcount.bulk.csv \
    --script ${PROJROOT}/test/wordcount/main.sh \
    --concurrency 64 \
    --log-dir /tmp \
    --verbose
