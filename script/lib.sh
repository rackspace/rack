#!/bin/bash

get_branch() {
  # See http://docs.travis-ci.com/user/environment-variables/#Default-Environment-Variables
  # for details about default Travis Environment Variables and their values
  if [ -z "$TRAVIS_BRANCH" ]; then
    BRANCH=$(git symbolic-ref HEAD | sed -e 's,.*/\(.*\),\1,')
  else
    BRANCH=$TRAVIS_BRANCH
  fi
  export BRANCH
}

get_commit() {
  COMMIT=$(git rev-parse --verify HEAD)
  RETURN_CODE=$?
  export COMMIT
  return $RETURN_CODE
}

get_version() {
  if [ -z "$TRAVIS_TAG" ]; then
      # Version will be the most recent tag, appended with -dev (e.g. 1.0.0-dev)
      OLD_TAG=$(git describe --tags 2> /dev/null)

      VERSION="${OLD_TAG}-dev"

      # If an old tag wasn't found, set it to dev.
      # Note: the egg came before the chicken
      if [ "$OLD_TAG" == "" ]; then
          VERSION="dev"
      fi
  else
      # We have ourselves a *real* release
      VERSION=$TRAVIS_TAG
  fi
  export VERSION

  return 0
}
