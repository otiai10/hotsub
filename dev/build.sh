#!/bin/sh

# This script compiles "awsub" binaries for multiple platforms,
# by using "gox" (https://github.com/mitchellh/gox).
# If you want to cross-compile 'awsub" binary by yourself,
# you need to run `go get -u github.com/mitchellh/gox" at first.

DEVELOPDIR=`cd $(dirname $0) && pwd`
PROJECTROOT=`dirname ${DEVELOPDIR}`
BUILDDIR=${PROJECTROOT}/builds

TARGETS="linux/amd64 darwin/amd64 windows/amd64"

rm -rf ${BUILDDIR}
mkdir -p ${BUILDDIR}

gox -output="${BUILDDIR}/{{.OS}}_{{.Arch}}/{{.Dir}}" -osarch="${TARGETS}" -rebuild -verbose

for dir in ${BUILDDIR}/*; do
    osarch=`basename ${dir}`
    tar -czvf builds/awsub.${osarch}.tar.gz -C ${dir} `ls ${dir}`
    rm -rf ${dir}
done
