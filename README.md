# rack
An opinionated CLI for Rackspace interactions

## Warning

This repo is under active development and is not ready for anything but testing
and development.

## Installation

Make sure you have Go installed and the GOPATH environment variable set.
```sh
go get github.com/jrperritt/rack
go build -o $GOPATH/bin/rack
```

Export the following environment variables:  
RS_REGION_NAME (DFW, IAD, ORD, LON, SYD, HKG)  
RS_USERNAME (Your Rackspace username)  
RS_AUTH_URL (https://identity.api.rackspacecloud.com/v2.0)  
RS_API_KEY (Your Rackspace API key)  

You should then be able to run commands:
```sh
rack compute servers list
```

## Bash Completion
Add the following line to your `.bashrc` file:
```sh
PROG=rack source $GOPATH/src/github.com/codegangsta/cli/autocomplete/bash_autocomplete
```
and source it:
```sh
source ~/.bashrc
```
