#!/bin/bash
set -x
# Check if the first command line argument is set
if [ -z "$1" ]
then
    echo "No arguments supplied. must be webhook-docker"
    exit 1
fi
echo "BUILDING $1"
BINARY_NAME=$1
BUILD_TARGET=cmd/${BINARY_NAME}/main.go

VERSION=$( git describe --tags --always --dirty --match=* 2> /dev/null || echo "v0.0.0")
SHORT_VERSION=$(git describe --tags --always --dirty --match=* | cut -d'-' -f1 2> /dev/null || echo "v0.0.0")
COMMIT=$(git rev-parse HEAD)
DATE=$( date +%FT%T%z)
LDFLAGS="-w -s -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.branch=${BRANCH} -X main.buildDate=${DATE} -X main.shortVersion=${SHORT_VERSION}"





echo "VERSION: $VERSION"
echo "SHORT_VERSION: $SHORT_VERSION"
echo "COMMIT: $COMMIT"
echo "BUILD DATE: $DATE"
echo "GOLANG FLAGS: $LDFLAGS"


echo "Build executable: go build -ldflags "${LDFLAGS}" -o ${BINARY_NAME} ${BUILD_TARGET}"
time  go build -ldflags "${LDFLAGS}" -o ${BINARY_NAME} ${BUILD_TARGET}

strip ${BINARY_NAME}
