#!/bin/bash

mkdir -p /var/in

if [ -n ${SOURCE} ]; then
  aws s3 cp ${SOURCE} /var/in
elif [ -n ${SOURCE_RECURSIVE} ]; then
  DIR=`basename ${SOURCE_RECURSIVE}`
  aws s3 cp -r ${SOURCE_RECURSIVE} /var/in/${DIR}
elif
  echo "[lifecycle error] Either of SOURCE or SOURCE_RECURSIVE must be specified" >&2
  exit 1
fi
