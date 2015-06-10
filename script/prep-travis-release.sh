#!/usr/bin/env bash

set -euo
BRANCH=$(git symbolic-ref HEAD | sed -e 's,.*/\(.*\),\1,')

mkdir -p build

BASENAME="rack"
BINARY=$BASENAME

# Only append the branch name on non-master branches
if [ $BRANCH != "master" ]; then
  BINARY=${BASENAME}-${BRANCH}
fi

# Build a version for this branch
go build -o build/$BINARY

# Provide a commit hash version
COMMIT=`git rev-parse HEAD 2> /dev/null`
cp build/$BINARY build/${BASENAME}-${COMMIT}
