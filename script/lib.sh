#!/bin/bash

# This script is intended to be `source`d

# It populates several environment variables via these functions:
# $BRANCH - get_branch()
# $VERSION - get_version()
# $COMMIT - get_commit()

get_branch() {
  # See http://docs.travis-ci.com/user/environment-variables/#Default-Environment-Variables
  # for details about default Travis Environment Variables and their values
  if [ -z "${TRAVIS_BRANCH-}" ]; then
    BRANCH=$(git symbolic-ref HEAD | sed -e 's,.*/\(.*\),\1,')
  else
    BRANCH=${TRAVIS_BRANCH}
  fi
  export BRANCH

  return 0
}

get_commit() {
  COMMIT=$(git rev-parse --verify HEAD)
  RETURN_CODE=$?
  export COMMIT
  return $RETURN_CODE
}

get_version() {
  if [ -z "${TRAVIS_TAG-}" ]; then
      # Version will be the most recent tag, appended with -dev (e.g. 1.0.0-dev)
      OLD_TAG=$(git describe --tags 2> /dev/null)

      VERSION="${OLD_TAG}-dev"

      # If an old tag wasn't found, set it to dev.
      # Note: the egg came before the chicken
      if [ "${OLD_TAG}" == "" ]; then
          VERSION="dev"
      fi
  else
      # We have ourselves a *real* release
      VERSION=${TRAVIS_TAG}
  fi
  export VERSION

  return 0
}

#
# Helper function to do replace; this should work across operating systems
#
update() {
  TMP_FILE=$(mktemp "$1")
  sed -e "$2" "$3" > "$TMP_FILE"
  chmod 0644 "$TMP_FILE"
  mv -f "$TMP_FILE" "$3"
}

update_docs() {
  NEW_VERSION=$1
  DOCS_INDEX_FILE="docs/index.rst"
  DOCS_CONFIGURATION_FILE="docs/configuration.rst"

  #
  # Update the docs index paths
  #

  update ./index.rst-tmpXXX "s#rackcdn\.com/[0-9a-zA-Z.-]*/#rackcdn\.com/$NEW_VERSION/#g" $DOCS_INDEX_FILE

  #
  # Update the docs configuration paths
  #

  update ./configuration.rst-tmpXXX "s#rackcdn\.com/[0-9a-zA-Z.-]*/#rackcdn\.com/$NEW_VERSION/#g" $DOCS_CONFIGURATION_FILE

}
