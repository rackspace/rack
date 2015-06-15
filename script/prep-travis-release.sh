#!/usr/bin/env bash

# prep-travis-release is used to layout a build directory for Travis to use for
# uploading to cloudfiles
#
# Assumes a gimme environment with GIMME_OS, GIMME_ARCH, and GOROOT, and PATH
#
# The script does not have to run on Travis though, as it will detect your
# current branch instead of using $TRAVIS_BRANCH

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

# Ensure GOARM is defined later
if [ "$arch" == "arm" -a -z "$GOARM" ]; then
  GOARM="6"
fi

# Back to strict
set -u
################################################################################

SUFFIX=""
if [ "$arch" == "arm" ]; then
  arch="armv${GOARM}"
fi
if [ "$os" == "windows" ]; then
  SUFFIX="${SUFFIX}.exe"
fi

echo "Building for ${os}-${arch}"

mkdir -p build
# Provide a commit hash version
COMMIT=`git rev-parse HEAD 2> /dev/null`
mkdir -p "build/commits/${COMMIT}"

BASENAME="rack-${os}-${arch}"
RACKBUILD="build/${BASENAME}"

go build -o $RACKBUILD

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
