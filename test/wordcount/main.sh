#!/bin/bash

echo "=== Parameter Check ==="
echo "HOTSUB_ROOT:  ${HOTSUB_ROOT}"
echo "SPEAKER_NAME: ${SPEAKER_NAME}"
echo "SPEECH_FILE:  ${SPEECH_FILE}"
echo "OUTDIR:       ${OUTDIR}"
echo "SLEEP:        ${SLEEP}"
echo "======================="

echo "=== Test --input-recursive ==="
cat ${META}/profile.txt
echo "======================="

echo "=== Common ENV ==="
echo "FOO: ${FOO}"
echo "=================="

cat ${SPEECH_FILE} \
  | tr ' ' '\n' \
  | tr -d , \
  | sort \
  | uniq -c \
  | sort -r \
  > ${OUTDIR}/word-count.txt
  # | tee ${OUTDIR}/word-count.txt

if [[ ${SLEEP} -gt 1 ]]; then
  echo "This script gonna sleep for ${SLEEP} seconds!"
  sleep ${SLEEP}s
fi

echo "Everything done on this machine"