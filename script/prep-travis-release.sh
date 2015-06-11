#!/usr/bin/env bash

set -eo

mkdir -p build
mkdir -p build/pr/

# Provide a commit hash version
COMMIT=`git rev-parse HEAD 2> /dev/null`
mkdir -p "build/commits/${COMMIT}"

echo '
  darwin   amd64    OS X
  freebsd  386      FreeBSD 32-bit
  freebsd  amd64    FreeBSD 64-bit
  linux    386      Linux 32-bit
  linux    amd64    Linux 64-bit
  windows  386      Windows 32-bit
  windows  amd64    Windows 64-bit
' | {
  while read os arch label; do
    [ -n "$os" ] || continue
    export GIMME_OS="$os"
    export GIMME_ARCH="$arch"

    BASENAME="rack-${os}-${arch}"
    RACKBUILD="build/${BASENAME}"

    SUFFIX=""
    if [ "$os" == "windows" ]; then
      SUFFIX=".exe"
    fi

    # Create build/rack
    go build -o $RACKBUILD

    # See http://docs.travis-ci.com/user/environment-variables/#Default-Environment-Variables
    # for details about default Travis Environment Variables and their values
    if [ -z "$TRAVIS_BRANCH" ]; then
      BRANCH=$(git symbolic-ref HEAD | sed -e 's,.*/\(.*\),\1,')
      cp $RACKBUILD build/${BASENAME}-${BRANCH}${SUFFIX}
    else
      # Ship a PR binary
      if [ "$TRAVIS_PULL_REQUEST" != "false" ]; then
        cp $RACKBUILD build/pr/${BASENAME}-pr${TRAVIS_PULL_REQUEST}${SUFFIX}
      fi

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

  done
}
