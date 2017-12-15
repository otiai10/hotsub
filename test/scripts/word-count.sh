#!/bin/bash

echo "=== Parameter Check ==="
echo "SPEECH_FILE: ${SPEECH_FILE}"
echo "OUTDIR:      ${OUTDIR}"
echo "======================="

cat ${SPEECH_FILE} \
  | tr ' ' '\n' \
  | tr -d , \
  | sort \
  | uniq -c \
  | sort -r \
  > ${OUTDIR}/wordcount.txt
