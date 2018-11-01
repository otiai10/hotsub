#!/bin/sh

PROJROOT=$(dirname $(dirname $(cd $(dirname $0) && pwd)))

set -e -v
hotsub run \
  --wdl ${PROJROOT}/test/wdl/hello.wdl \
  --wdl-job ${PROJROOT}/test/wdl/job-0.json \
  --wdl-job ${PROJROOT}/test/wdl/job-1.json \
  --verbose

  # --provider gcp \
  # --google-project genomondevel1 \
