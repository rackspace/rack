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

declare -xr CDN="https://ec4a542dbf90c03b9f75-b342aba65414ad802720b41e8159cf45.ssl.cf5.rackcdn.com"
declare -xr BUILDDIR="build"


################################################################################
# Disable strict temporarily to accept global environment variables that come
# from GIMME and Travis
################################################################################
set +u

if [[ -z "$GIMME_OS" && -z "$GIMME_ARCH" ]]; then
  >&2 echo "GIMME_OS and GIMME_ARCH must be defined"
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
    >&2 echo "Unknown OS ${os}. Assuming it's a valid OS for gimme/go and charging ahead."
esac

case $arch in
  arm)
    # Note that we can never be 1:1 between `uname -m` and arm versions
    arch="armv${GOARM}"
    ;;
esac

SUFFIX=""
if [ "$os" == "Windows" ]; then
  SUFFIX=".exe"
fi

################################################################################
# Set up the build and deploy layout
################################################################################

# Get the version from travis, or default to 0.0.0
if [ -z $TRAVIS_TAG ]; then
  VERSION=$TRAVIS_TAG
else
  VERSION="v0.0.0"
fi

BASEDIR="${VERSION}/${os}/${arch}"
# Mirror the github layout for branches, tags, commits
TREEDIR="${os}/${arch}/tree"

mkdir -p $BUILDDIR
mkdir -p $BUILDDIR/$BASEDIR
mkdir -p $BUILDDIR/$TREEDIR

BASENAME="rack"

# Base build not in build dir to prevent accidental upload on failure
RACKBUILD="${BASENAME}${SUFFIX}"

COMMIT=$(git rev-parse --verify HEAD)
sed -i "s/var Commit =.*/var Commit = \"$COMMIT\"/" util/commit.go

go build -o $RACKBUILD

# Ship /tree/rack-branchname
cp $RACKBUILD ${BUILDDIR}/${TREEDIR}/${BASENAME}-${BRANCH}${SUFFIX}
echo "Fresh build for branch '${BRANCH}' at "
echo "${CDN}/${TREEDIR}/${BASENAME}-${BRANCH}${SUFFIX}"

if [ "$BRANCH" == "master" ]; then
  # Only when we're on master do we spit out the official ones.
  cp $RACKBUILD ${BUILDDIR}/${BASEDIR}/${BASENAME}${SUFFIX}
  echo "Get it while it's hot at"
  echo "${CDN}/${BASEDIR}/${BASENAME}${SUFFIX}"
fi

# Clean up after build
rm $RACKBUILD
