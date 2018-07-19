#!/bin/sh

PROJROOT=$(dirname $(dirname $(cd $(dirname $0) && pwd)))

go install .

set -e -v
hotsub \
  --cwl ${PROJROOT}/test/cwl/hello.cwl \
  --cwl-param ${PROJROOT}/test/cwl/job.yml \
  --verbose