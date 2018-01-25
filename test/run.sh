#!/bin/bash
# set -e

PROJROOT=$(dirname $(cd $(dirname $0) && pwd))

for script in ${PROJROOT}/test/cases/*.sh; do

  # Decide expected exit status
  name=`basename ${script}`
  if [[ ${name} = *".FAIL.sh" ]]; then expected=1; else expected=0; fi

  # Execute testcase script
  ${script}

  # Assert exit status
  if [[ $? != ${expected} ]]; then
    echo "NG! ${name}"
    exit 1
  else
    # Pass
    echo "ok ${name}"
  fi

done
