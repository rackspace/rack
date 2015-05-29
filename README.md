# rackcli
An opinionated CLI for Rackspace interactions

## Warning

This repo is under active development and is not ready for anything but testing
and development.

## Installation

Make sure you have Go installed and the GOPATH environment variable set.
```sh
go get github.com/jrperritt/rackcli
go build -o $GOPATH/bin/rackcli
```

Export the following environment variables:
RS_REGION_NAME
RS_USERNAME
RS_AUTH_URL
RS_API_KEY

You should then be able to run commands:
```sh
rackcli compute servers list
```
