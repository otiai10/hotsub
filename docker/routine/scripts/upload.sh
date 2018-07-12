#!/bin/bash

# Let it fail when something failed
set -e

# This script is supposed to be used to upload output files to
# cloud storage services such as Amazon S3, GoogleStorage and Azure Storage.
# This script should be called for just 1 time to upload "SRC" directory.
#
# Supported Services are:
#
# - Amazon S3
# - Google Cloud Storage
#
# it means the docker image called "hotsub/routine" must have
# following binaries in its PATH:
#
# - aws (aws s3)
# - gsutil
#
# Example:
#
#     SOURCE=/your/outputs DEST=s3://your-bucket/your-object ./upload.sh

function precheck() {
  if [[ -z ${SOURCE} ]]; then
    echo "SOURCE variable is required but not provided"
    exit 1
  fi
  if [[ -z ${DEST} ]]; then
    echo "DEST variable is required but not provided"
    exit 1
  fi
}

function upload() {


  if [[ ${DEST} =~ s3://.+ ]]; then
    PROVIDER=s3
  elif [[ ${DEST} =~ gs://.+ ]]; then
    PROVIDER=gs
  else
    echo "Unknown destination format: ${DEST}"
    exit 1
  fi

  case ${PROVIDER} in
  s3)
    CMD="aws s3 cp --recursive"
    ;;
  gs)
    CMD="gsutil cp -R"
    ;;
  esac

  ${CMD} ${SOURCE} ${DEST} || exit $?
}

function __main__() {
  precheck
  upload
}

__main__ $@
