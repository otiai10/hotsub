#!/bin/bash

# Let it fail when something failed
set -e

# This script is supposed to be used to download input files from
# cloud storage services such as Amazon S3, GoogleStorage and Azure Storage.
# This script should be called for each 1 input (or input-recursive) file.
#
# Supported Services are:
#
# - Amazon S3
# - Google Cloud Storage
#
# it means the docker image called "awsub/lifecycle" must have
# following binaries in its PATH:
#
# - aws (aws s3)
# - gsutil
#
# Example:
#
#     INPUT=s3://your-bucket/your-object DIR=/your/path ./download.sh

function precheck() {
  if [[ -z ${INPUT} && -z ${INPUT_RECURSIVE} ]]; then
    echo "Either of INPUT or INPUT_RECURSIVE must be specified to invoke 'download.sh'"
    exit 1
  fi
}

function download() {

  mkdir -p ${DIR}

  if [[ -n ${INPUT} ]]; then SRC=${INPUT}; elif [[ -n ${INPUT_RECURSIVE} ]]; then SRC=${INPUT_RECURSIVE}; fi
  if [[ ${SRC} =~ s3://.+ ]]; then
    PROVIDER=s3
  elif [[ ${SRC} =~ gs://.+ ]]; then
    PROVIDER=gs
  else
    echo "Provided input doesn't have '{provider}://{your-bucket}' format"
    exit 1
  fi

  case ${PROVIDER} in
  s3)
    if [[ -n ${INPUT_RECURSIVE} ]]; then
      # CMD="aws s3 cp --recursive"
      CMD="aws s3 sync"
    else
      CMD="aws s3 cp"
    fi
    ;;
  gs)
    if [[ -n ${INPUT_RECURSIVE} ]]; then
      CMD="gsutil cp -R"
    else
      CMD="gsutil cp"
    fi
    ;;
  esac

  DEST=${DIR}/`basename ${SRC}`
  echo "Execution: ${CMD} ${SRC} ${DEST}"
  ${CMD} ${SRC} ${DEST} || exit $?
}

function __main__() {
  precheck
  download
}

__main__ $@
