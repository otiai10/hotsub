#!/bin/bash

echo "=== Parameter Check ==="
echo "USER:         ${USER}"
echo "SPEAKER_NAME: ${SPEAKER_NAME}"
echo "SPEECH_FILE:  ${SPEECH_FILE}"
echo "DIR:          ${DIR}"
echo "======================="

cat ${SPEECH_FILE} \
  | tr ' ' '\n' \
  | tr -d , \
  | sort \
  | uniq -c \
  | sort -r \
# > ${DIR}/wordcount.txt
