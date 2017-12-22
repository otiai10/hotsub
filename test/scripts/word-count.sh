#!/bin/bash

echo "=== Parameter Check ==="
echo "USER:         ${USER}"
echo "SPEAKER_NAME: ${SPEAKER_NAME}"
echo "SPEECH_FILE:  ${SPEECH_FILE}"
echo "DIR:          ${DIR}"
echo "OUTDIR:       ${OUTDIR}"
echo "======================="

cat ${SPEECH_FILE} \
  | tr ' ' '\n' \
  | tr -d , \
  | sort \
  | uniq -c \
  | sort -r \
  | tee ${OUTDIR}/word-count.txt
