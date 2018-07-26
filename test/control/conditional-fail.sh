#!/bin/bash

if [ -n "${FAIL}" ]; then
  echo >&2 "Let it fail on purpose."
  exit 1
fi

count="0"
while [ $count -lt 120 ]; do
  count=$((count + 1))
  sleep 0.5s
  printf "."
done

echo >&1 "Finished"