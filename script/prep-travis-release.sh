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

if [[ -z "$GIMME_OS" && -z "$GIMME_ARCH" ]]; then
  echo "GIMME_OS and GIMME_ARCH must be defined"
  exit 2
fi

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

case $os in
  windows)
    os="Windows"
    ;;
  linux)
    os="Linux"
    ;;
  darwin)
    os="Darwin"
    ;;
  freebsd)
    os="FreeBSD"
    ;;
  *)
    echo "Unknown OS ${os}. Assuming it's a valid OS for gimme/go and charging ahead."
esac

case $arch in
  arm)
    # Note that we can never be 1:1 between `uname -m` and arm versions
    arch="armv${GOARM}"
    SUFFIX="-${os}-${arch}"
    ;;
esac

if [[ "$os" == "Windows" && "$arch" == "amd64" ]]; then
  SUFFIX=".exe"
elif [[ "$os" == "Windows" && "$arch" != "amd64" ]]; then
  SUFFIX="-${arch}.exe"
elif [ "$arch" == "amd64" ]; then
  # Assume 64 bit gets to be as is
  SUFFIX="-${os}"
else
  SUFFIX="-${os}-${arch}"
fi

echo "Building for ${os}-${arch}"

mkdir -p build
# Provide a commit hash version
COMMIT=`git rev-parse HEAD 2> /dev/null`
mkdir -p "build/commits/${COMMIT}"

BASENAME="rack"
RACKBUILD="build/${BASENAME}${SUFFIX}"

go build -o $RACKBUILD

cp $RACKBUILD build/commits/${COMMIT}/${BASENAME}-${COMMIT}${SUFFIX}

if [ "$BRANCH" != "master" ]; then
  # Ship /rack-branchname
  cp $RACKBUILD build/${BASENAME}-${BRANCH}${SUFFIX}
  # Remove our artifact
  rm $RACKBUILD
fi
