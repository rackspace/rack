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
go get github.com/rackspace/rack
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
`rack --username user1 --api-key 123456789 --region DFW servers instance list`

### Config file
You can create a config file in ` ~/.rack/config`:

```
[DEFAULT]
username=user1
api-key=123456789
region=DFW

[PROFILE2]
username=user2
api-key=987654321
region=IAD
```

If you're using the default profile, you can call a command without additional flags:
```sh
rack servers instance list
```

However, if you'd like to use a different profile (such as PROFILE2 above):
```sh
rack --profile PROFILE2 servers instance list
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

## Roadmap

Our plan is to release a limited beta of `rack` the first week of August 2015 with full support for the following services:

- Cloud Servers
- Cloud Files
- Cloud Networks
- Cloud Block Storage

We expect to exit beta in the coming weeks, while also focusing on adding new services. The current plan of record to add new services is based on the following priority:

1. Cloud Load Balancers
2. Cloud DNS \*
3. Cloud Databases \*
4. Cloud Identity
5. Cloud Images \*
6. Rackspace CDN
7. RackConnect
8. Cloud Big Data \*
9. Cloud Monitoring \*
10. Cloud Orchestration
11. Cloud Queues \*
12. Cloud Backup \*
13. Autoscale \*
14. Cloud Metrics \*

\* Services not supported in [Gophercloud](github.com/rackspace/gophercloud) at present.
