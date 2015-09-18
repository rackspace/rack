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
declare -xr BASENAME="rack"

SCRIPT_DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
source "${SCRIPT_DIR}/lib.sh"

initialize() {
  if [[ -z "${GIMME_OS-}" && -z "${GIMME_ARCH-}" ]]; then
    >&2 echo "GIMME_OS and GIMME_ARCH must be defined"
    exit 2
  fi

  os=$GIMME_OS
  arch=$GIMME_ARCH

  get_branch
  get_version
  get_commit

  # Ensure GOARM is defined later
  if [ "$arch" == "arm" -a -z "${GOARM-}" ]; then
    GOARM="6"
  fi

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
}

build() {
  RACKBUILD="${BASENAME}${SUFFIX}"
  script/build ${RACKBUILD}
}

deploy() {
  # Set up the directory structure for release to rack's backing CDN
  # and copy the built `rack` into the right paths for this build type

  BASEDIR="${VERSION}/${os}/${arch}"
  # Mirror the github layout for branches, tags, commits
  TREEDIR="${os}/${arch}/tree"

  mkdir -p "${BUILDDIR}"
  mkdir -p "${BUILDDIR}/${BASEDIR}"
  mkdir -p "${BUILDDIR}/${TREEDIR}"

  if (( $? != 0 )); then
    echo "Failed build."
    exit 1
  fi

  # Ship /tree/rack-branchname
  cp "${RACKBUILD}" "${BUILDDIR}/${TREEDIR}/${BASENAME}-${BRANCH}${SUFFIX}"

  # Only when we're on the canonical rackspace/rack repo will we be shipping
  # binaries, which comes down to whether TRAVIS_SECURE_ENV_VARS is defined
  if [ -n "${TRAVIS_SECURE_ENV_VARS-}" ]; then
    echo "Fresh build for branch '${BRANCH}' at "
    echo "${CDN}/${TREEDIR}/${BASENAME}-${BRANCH}${SUFFIX}"
    if [ -n "${TRAVIS_TAG-}" ]; then
      # Only when we're on an official tag do we spit out the official ones.
      cp "${RACKBUILD}" "${BUILDDIR}/${BASEDIR}/${BASENAME}${SUFFIX}"
      echo "Get it while it's hot at"
      echo "${CDN}/${BASEDIR}/${BASENAME}${SUFFIX}"
    fi
    # Clean up after build
    rm $RACKBUILD
  else
    # Do nothing, keep the built artifact
    echo "${RACKBUILD} is ready for you"
  fi
}

initialize
build
deploy
