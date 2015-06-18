# rack
An opinionated CLI for Rackspace interactions

## Warning

This repo is under active development and is not ready for anything but testing
and development.

## Documentation

For complete documentation, see [the docs](http://rackspace-cli.readthedocs.org/en/latest/).

## Installation

Make sure you have Go installed and the GOPATH environment variable set.
```sh
go get github.com/jrperritt/rack
go build -o $GOPATH/bin/rack
```

## Setting Authentication Credentials

### Environment Variables
Export the following environment variables:  
`RS_REGION_NAME` (DFW, IAD, ORD, LON, SYD, HKG)  
`RS_USERNAME` (Your Rackspace username)  
`RS_API_KEY` (Your Rackspace API key)  

### Command-line
You can set auth parameters on the command-line using global flags:
`rack --username user1 --apikey 123456789 --region DFW servers instance list`

### Config file
You can create a config file in ` ~/.rack/config`:

```
[DEFAULT]
username=user1
apikey=123456789
region=DFW

[PROFILE2]
username=user2
apikey=987654321
region=IAD
```

If you're using the default profile, you can call a command without additional flags:
```sh
rack compute servers list
```

However, if you'd like to use a different profile (such as PROFILE2 above):
```sh
rack --profile PROFILE2 compute servers list
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
