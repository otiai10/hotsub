#!/bin/sh

PROJROOT=$(dirname $(dirname $(cd $(dirname $0) && pwd)))

set -e -v
hotsub \
  --cwl ${PROJROOT}/test/cwl/hello.cwl \
  --cwl-param ${PROJROOT}/test/cwl/job.yml \
  --verbose