#!/bin/bash

set -e -u

DEST=${OUTDIR}/${CASENAME}

bwa mem \
    -t 2 -T 0 \
    ${REFRENCE}/GRCh37.fa \
    ${FASTQ_1} ${FASTQ_2} \
    > ${DEST}/result.sam

samtools view \
    -S -b -h \
    ${DEST}/result.sam \
    > ${DEST}/result.bam

samtools sort \
    ${DEST}/result.bam \
    > ${DEST}/sorted.bam