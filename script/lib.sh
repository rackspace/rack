get_branch() {
  # See http://docs.travis-ci.com/user/environment-variables/#Default-Environment-Variables
  # for details about default Travis Environment Variables and their values
  if [ -z "$TRAVIS_BRANCH" ]; then
    BRANCH=$(git symbolic-ref HEAD | sed -e 's,.*/\(.*\),\1,')
  else
    BRANCH=$TRAVIS_BRANCH
  fi
}

get_commit() {
  COMMIT=$(git rev-parse --verify HEAD)
  return $?
}

get_version() {
  if [ -z "$TRAVIS_TAG" ]; then
      # Version will be the most recent tag, appended with -dev (e.g. 1.0.0-dev)
      OLD_TAG=$(git describe --tags 2> /dev/null)
      VERSION="${OLD_TAG}-dev"
      if [ "$OLD_TAG" == "" ]; then
          VERSION="dev"
      fi
  else
      # We have ourselves a *real* release
      VERSION=$TRAVIS_TAG
  fi

  return 0
}
