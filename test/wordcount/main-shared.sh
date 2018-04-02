#!/bin/bash

set -e -u

echo "=== Parameter Check ==="
echo "AWSUB_ROOT:   ${AWSUB_ROOT}"
echo "SPEAKER_NAME: ${SPEAKER_NAME}"
echo "SPEECH_FILE:  ${SPEECH_FILE}"
echo "OUTDIR:       ${OUTDIR}"
echo "======================="

echo "=== Test --input-recursive ==="
cat ${META}/profile.txt
echo "======================="

echo "===== Test --shared ====="
cat ${NYANCAT}/nyancat.txt
echo "========================="

cat ${SPEECH_FILE} \
  | tr ' ' '\n' \
  | tr -d , \
  | sort \
  | uniq -c \
  | sort -r \
  > ${OUTDIR}/word-count.txt
  # | tee ${OUTDIR}/word-count.txt
