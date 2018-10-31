#!/bin/bash

PROJROOT=$(dirname $(dirname $(cd $(dirname $0) && pwd)))

set -e -v
hotsub run \
    --tasks ${PROJROOT}/test/wordcount/wordcount.gcp.csv \
    --script ${PROJROOT}/test/wordcount/main-shared.sh \
    --shared NYANCAT=gs://hotsub/speech/shared \
    --shareddata-disksize 16 \
    --provider gcp \
    --google-zone asia-northeast1-a \
    --google-project genomondevel1 \
    --log-dir /tmp \
    --verbose
