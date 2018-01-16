#!/bin/sh

set -o errexit
set -o nounset
set -o xtrace

STAR --version

OUTPUT_PREF=${OUTPUT_DIR}/${SAMPLE}

STAR_OPTION="--runThreadN 6 --outSAMstrandField intronMotif --outSAMunmapped Within --outSAMtype BAM Unsorted"
SAMTOOLS_SORT_OPTION="-@ 6 -m 3G"

ls -lh ${REFERENCE}

/usr/local/bin/STAR \
    --genomeDir ${REFERENCE} \
    --readFilesIn ${INPUT1} ${INPUT2} \
    --outFileNamePrefix ${OUTPUT_PREF}. \
    ${STAR_OPTION}

/usr/local/bin/samtools sort \
    -T ${OUTPUT_PREF}.Aligned.sortedByCoord.out \
    ${SAMTOOLS_SORT_OPTION} \
    ${OUTPUT_PREF}.Aligned.out.bam -O bam \
    > ${OUTPUT_PREF}.Aligned.sortedByCoord.out.bam

/usr/local/bin/samtools index \
    ${OUTPUT_PREF}.Aligned.sortedByCoord.out.bam

rm ${OUTPUT_PREF}.Aligned.out.bam

echo "Finished ==> ${SAMPLE}"
