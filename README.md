# rack
The official command line interface for the Rackspace Cloud.

`rack` provides a consistent interface for interacting with the Rackspace Cloud. For example, creating a new server:

```
$ rack servers instance create \
    --name my-server \
    --image-id 5ed162cc-b4eb-4371-b24a-a0ae73376c73 \
    --flavor-id general1-1 \
    --keypair my-laptop

ID		9818861f-2f14-437f-89b0-a36dfa1831b7
AdminPass	4vLb2PiqUGdP
```

For complete documentation, see [the docs](https://developer.rackspace.com/docs/rack-cli/).

## Download

`rack` binaries are available for OSX, Linux, and Windows. [Download `rack` now](https://developer.rackspace.com/docs/rack-cli/configuration#installation).

## Build From Source

Make sure you have [Go installed](https://golang.org/doc/install) and the [`GOPATH`](https://golang.org/doc/code.html#GOPATH) environment variable set.
```sh
go get github.com/rackspace/rack
cd $GOPATH/src/github.com/rackspace/rack
./script/build
```

## Setting Authentication Credentials

`rack` has a number of ways of getting credentials configured. The easiest way is to use `rack configure` to setup a profile:

```
$ rack configure

    This interactive session will walk you through creating
    a profile in your configuration file. You may fill in all or none of the
    values.

    Rackspace Username: iamacat
    Rackspace API key: secrets
    Rackspace Region: IAD
    Profile Name (leave blank to create a default profile):
```

This will create a configuration file at `~/.rack/config` which will store your credentials. When using the default profile, you don't need to specify the profile:

```
rack servers instance list
```

However, if you've named a profile, you can use the `--profile` flag to specify it explicitly:

```
rack servers instance list --profile staging
```

### Environment Variables

In addition to using the config file, you can also use environment variables. The following environment variables are supported:

`RS_REGION_NAME` (DFW, IAD, ORD, LON, SYD, HKG)  
`RS_USERNAME` (Your Rackspace username)  
`RS_API_KEY` (Your Rackspace API key)  

### Command-line

Lastly, you can also set your authentication credentials as flags:

```

$ rack servers instance list --username user1 --api-key 123456789 --region DFW

```

### Precedence

`rack` will first look for credentials as flags. Secondly, `rack` will use a named profile, or `default` if there is a default profile. Lastly, `rack` will look for environment variables.

## Bash Completion

At any time, you can run `rack init` which will create the auto-completion file in `~/.rack/bash_autocomplete` and add it to `~/.bash_profile`. You'll need to restart your terminal session to enable auto-completion.

For command completion on Windows, please see the [command completion](https://developer.rackspace.com/docs/rack-cli/configuration#command-completion) section on the configuration page of our documentation.

## Roadmap

Below is a tentative list of services to be added. Though currently ordered by priority, the list can and will change:

1. Cloud Load Balancers
2. Cloud DNS \*
3. Cloud Databases \*
4. Cloud Identity
5. Cloud Images \*
6. Rackspace CDN
7. RackConnect
8. Cloud Big Data \*
9. Cloud Monitoring \*
10. [~~Cloud Orchestration~~](https://github.com/rackspace/rack/commit/5d311719579c98a609beb7801c593e132711e003)
11. Cloud Queues \*
12. Cloud Backup \*
13. Autoscale \*
14. Cloud Metrics \*

\* Services not supported in [Gophercloud](https://github.com/rackspace/gophercloud) at present.
