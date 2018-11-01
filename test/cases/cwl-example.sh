#!/bin/sh

PROJROOT=$(dirname $(dirname $(cd $(dirname $0) && pwd)))

set -e -v
hotsub run \
  --cwl ${PROJROOT}/test/cwl/hello.cwl \
  --cwl-job ${PROJROOT}/test/cwl/job.yml \
  --verbose