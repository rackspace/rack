#!/usr/bin/env bash

set -eo

mkdir -p build
mkdir -p build/pr/

# Provide a commit hash version
COMMIT=`git rev-parse HEAD 2> /dev/null`
mkdir -p "build/commits/${COMMIT}"

eval "$(gimme 1.4)"

os=$GIMME_OS
arch=$GIMME_ARCH

echo "Pre-flight diagnostics"
ls
pwd

echo "Building for ${os}-${arch}"

BASENAME="rack-${os}-${arch}"
RACKBUILD="build/${BASENAME}"

SUFFIX=""
if [ "$os" == "windows" ]; then
  SUFFIX=".exe"
fi

# Create build/rack, making sure to start over without build artifacts
go build -a -o $RACKBUILD

# See http://docs.travis-ci.com/user/environment-variables/#Default-Environment-Variables
# for details about default Travis Environment Variables and their values
if [ -z "$TRAVIS_BRANCH" ]; then
  BRANCH=$(git symbolic-ref HEAD | sed -e 's,.*/\(.*\),\1,')
else
  BRANCH=$TRAVIS_BRANCH
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
