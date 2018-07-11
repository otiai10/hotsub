#!/bin/bash

PROJROOT=$(dirname $(dirname $(cd $(dirname $0) && pwd)))

set -e -v
hotsub \
    --tasks ${PROJROOT}/test/wordcount/wordcount.csv \
    --script ${PROJROOT}/test/wordcount/main-shared.sh \
    --shared NYANCAT=s3://hotsub/examples/wordcount/nyancat \
    --shareddata-disksize 8 \
    --aws-shared-instance-type t2.micro \
    --log-dir /tmp \
    --verbose
