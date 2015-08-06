#!/bin/bash
#
# Install the Sphinx deconst preparer and submit documentation to developer.rackspace.com.

set -euo pipefail

pip3 install -e git+https://github.com/deconst/preparer-sphinx.git#egg=deconstrst

export CONTENT_STORE_APIKEY=${APIKEY1}${APIKEY2}${APIKEY3}

cd docs/
deconst-preparer-sphinx
