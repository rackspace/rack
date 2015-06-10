#!/usr/bin/env bash

set -eo

mkdir -p build

BASENAME="rack"
SUFFIX=""

# See http://docs.travis-ci.com/user/environment-variables/#Default-Environment-Variables
# for details about default Travis Environment Variables and their values
if [ -z "$TRAVIS_BRANCH" ]; then
  BRANCH=$(git symbolic-ref HEAD | sed -e 's,.*/\(.*\),\1,')
  SUFFIX="-${BRANCH}"
else
  # TRAVIS PULL REQUEST is literally "false" when not a PR, number otherwise
  if [ -z "$TRAVIS_PULL_REQUEST" ] || [ "$TRAVIS_PULL_REQUEST" == "false" ]; then
    SUFFIX="-$TRAVIS_BRANCH"

    # No branch name for master
    if [ "master" == "$TRAVIS_BRANCH" ]; then
      SUFFIX=""
    fi

  else
    SUFFIX="-${TRAVIS_PULL_REQUEST}"
  fi
fi

BINARY="${BASENAME}${SUFFIX}"

# Build a version for this branch
go build -o build/$BINARY

# Provide a commit hash version
COMMIT=`git rev-parse HEAD 2> /dev/null`
cp build/$BINARY build/${BASENAME}-${COMMIT}
