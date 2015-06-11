#!/usr/bin/env bash

set -euo pipefail
IFS=$'\n\t'


################################################################################
# Disable strict temporarily to accept global environment variables that come
# from GIMME and Travis
################################################################################
set +u

os=$GIMME_OS
arch=$GIMME_ARCH
# See http://docs.travis-ci.com/user/environment-variables/#Default-Environment-Variables
# for details about default Travis Environment Variables and their values
if [ -z "$TRAVIS_BRANCH" ]; then
  BRANCH=$(git symbolic-ref HEAD | sed -e 's,.*/\(.*\),\1,')
else
  BRANCH=$TRAVIS_BRANCH
fi
# Back to strict
set -u
################################################################################

echo "Building for ${os}-${arch}"

mkdir -p build
# Provide a commit hash version
COMMIT=`git rev-parse HEAD 2> /dev/null`
mkdir -p "build/commits/${COMMIT}"

BASENAME="rack-${os}-${arch}"
RACKBUILD="build/${BASENAME}"

go build -o $RACKBUILD

SUFFIX=""
if [ "$os" == "windows" ]; then
  SUFFIX=".exe"
fi

cp $RACKBUILD build/commits/${COMMIT}/${BASENAME}-${COMMIT}${SUFFIX}

if [ "$BRANCH" != "master" ]; then
  # Ship /rack-branchname
  cp $RACKBUILD build/${BASENAME}-${BRANCH}${SUFFIX}
  # Remove our artifact
  rm $RACKBUILD
elif [ "$os" == "windows" ]; then
  cp $RACKBUILD $RACKBUILD${SUFFIX}
  rm $RACKBUILD
fi
