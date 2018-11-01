#!/bin/sh

PROJROOT=$(dirname $(dirname $(cd $(dirname $0) && pwd)))

set -e -v
hotsub run \
  --cwl ${PROJROOT}/test/cwl/js-docker.cwl \
  --cwl-job ${PROJROOT}/test/cwl/js-job.yml \
  --include ${PROJROOT}/test/cwl/hello.js \
  --verbose