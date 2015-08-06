# rack
The official command line interface for the Rackspace Cloud

## Warning

This repo is under active development and is not ready for anything but testing
and development.

## Documentation

For complete documentation, see [the docs](http://rackspace-cli.readthedocs.org/en/latest/).

## Download

Multi-platform binaries are available to download at [https://developer.rackspace.com/downloads](link tbd).

## Build from Source

Make sure you have [Go installed](https://golang.org/doc/install) and the [`GOPATH`](https://golang.org/doc/code.html#GOPATH) environment variable set.
```sh
go get github.com/rackspace/rack
go build -o $GOPATH/bin/rack
```

## Setting Authentication Credentials

`rack` has a number of ways of getting credentials configured. The easiest way is to use `rack configure` to setup a profile:

```bash
$rack configure

    This interactive session will walk you through creating
    a profile in your configuration file. You may fill in all or none of the
    values.

    Rackspace Username: iamacat
    Rackspace API key: secrets
    Rackspace Region : IAD
    Profile Name (leave blank to create a default profile):
```

This will create a configuration file at `~/.rack/config` which will store your credentials. When using the default profile, you don't need to specify the profile:

```bash
rack servers instance list
```

However, if you've named a profile, you can use the `--profile` flag to specify it explicitly:

```bash
rack servers instance list --profile staging
```

### Environment Variables

In addition to using the config file, you can also use environment variables. The following environment variables are supported:

`RS_REGION_NAME` (DFW, IAD, ORD, LON, SYD, HKG)  
`RS_USERNAME` (Your Rackspace username)  
`RS_API_KEY` (Your Rackspace API key)  

### Command-line

Lastly, you can also set your authentication credentials as flags:

```bash

$ rack servers instance list --username user1 --api-key 123456789 --region DFW

```

## Bash Completion

At any time, you can run `rack init` which will create the auto-completion file in `~/.rack/bash_autocomplete` and add it to `~/.bash_profile`. You'll need to restart your terminal session to enable auto-completion.

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

\* Services not supported in [Gophercloud](https://github.com/rackspace/gophercloud) at present.
