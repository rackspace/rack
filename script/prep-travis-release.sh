#!/usr/bin/env bash

set -eo

mkdir -p build
mkdir -p build/commits/
mkdir -p build/pr/

BASENAME="rack"
SUFFIX=""

RACKBUILD="build/${BASENAME}"

# Create build/rack
go build -o $RACKBUILD

# See http://docs.travis-ci.com/user/environment-variables/#Default-Environment-Variables
# for details about default Travis Environment Variables and their values
if [ -z "$TRAVIS_BRANCH" ]; then
  BRANCH=$(git symbolic-ref HEAD | sed -e 's,.*/\(.*\),\1,')
  cp $RACKBUILD build/${BASENAME}-${BRANCH}
else
  # Ship a PR binary
  if [ "$TRAVIS_PULL_REQUEST" != "false" ]; then
    cp $RACKBUILD build/pr/${BASENAME}-pr${TRAVIS_PULL_REQUEST}
  fi

  BRANCH=$TRAVIS_BRANCH
fi

# Ship /rack-branchname
cp $RACKBUILD build/${BASENAME}-${BRANCH}

# Provide a commit hash version
COMMIT=`git rev-parse HEAD 2> /dev/null`
cp $RACKBUILD build/commits/${BASENAME}-${COMMIT}
