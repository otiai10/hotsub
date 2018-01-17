#!/bin/sh

echo "[debug] Check uname"
uname -a

echo "[debug] Check tesseract version"
tesseract --version

# Execute OCR now
tesseract ${IMAGE} ${OUTDIR}/out

echo "[debug] Completed. Bye."
