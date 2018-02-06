#!/bin/bash
# set -e

PROJROOT=$(dirname $(cd $(dirname $0) && pwd))

MATCH=
while [[ $# -gt 0 ]]; do
case "${1}" in
  --match|-m)
  MATCH="${2}"
  shift && shift
  ;;
  *)
  shift
  ;;
esac
done

for script in ${PROJROOT}/test/cases/*.sh; do

  # Decide expected exit status
  name=`basename ${script}`

  if [[ -n ${MATCH} ]]; then
    if [[ ${name} != *${MATCH}* ]]; then continue; fi
  fi

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
