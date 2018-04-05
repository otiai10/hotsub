#!/bin/bash

set -e -u

DEST=${OUTDIR}/${CASENAME}
mkdir -p ${DEST}

set -v
bwa mem \
    -t 2 -T 0 \
    ${REFERENCE}/GRCh37.fa \
    ${FASTQ_1} ${FASTQ_2} \
    > ${DEST}/result.sam
samtools view \
    -S -b -h \
    ${DEST}/result.sam \
    > ${DEST}/result.bam
samtools sort \
    ${DEST}/result.bam \
    > ${DEST}/sorted.bam