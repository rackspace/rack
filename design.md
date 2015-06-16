# Rack CLI Design Document

This document will attempt to outline the design for the Rackspace Global Command Line Interface. It will not exhaustively define each interface; rather, it will focus on the philosophy and representative examples that should be used to inform the successive specifications for the specific service implementations.

As such, there are 4 core tenets we're using to define the design, in order of priority:

1. Correctness
2. Consistency
3. Completeness
4. Terseness

## Correctness

As a top level objective, the CLI must operate correctly. It's unacceptable to deliver a non-functional tool, nor one that behaves contrary to the documented expectations. The key here is to document the expected behavior, and then ensure that the delivered tool matches those expectations.

## Consistency

Our ambition for consistency is to provide a model that is used throughout the command line interface. All calls take the following format:

```
rack <service> [--options] <model> <action> [--flags]
```

An example of `service` might be `servers` or `files` or `networks`. `model` might be `instance`, `object` or `container`, or perhaps `network`. Action would be the underlying API call to make, for example `create` or `delete`.

Specifying `options` might be explicitly providing authentication information, for example `username` or `apiKey`, or other options that are not specific to the underlying `action` flags. The complete list of valid global options will be listed in a separate specification.

All of the following are examples of the pattern being applied consistently:

* `rack servers instance list` - Get a list of all instances
* `rack files container list` - Get a list of all containers
* `rack networks security-group get --id 12345` - Get a security group by id
* `rack load-balancers node remove --id 12345 --load-balancer-id abcdef` - Remove a node from a load balancer
* `rack files --username foobar --apiKey baz --region ORD object list --container mycontainer` - List all of the objects in container `mycontainer` in ORD for the user `foobar`

## Completeness

Our ultimate goal is to have complete coverage of every public feature across the breadth of the Rackspace Public Cloud.

Our initial releases may choose to ship a subset of services, but should be feature complete within those services. For example, we may choose for the initial release to support Cloud Servers and Cloud Networks, and Cloud Block Storage only, but these services should be feature complete. Subsequent releases should add one or more fully functional services. We should not release support for a service that is incomplete.

## Terseness

To the extent of the preceding objectives, terseness is the last priority. We should, so long as it does not come at the expense of the goals listed above, optimize the command line interface for efficiency.

For example, we may be able to imply values without the `--id` flag as appropriate:

```
rack servers instance get --id 12345
```

might be synonymous with:

```
rack servers instance get 12345
```

This is an example of a way to add terseness without compromising consistency. The key in this example would be that *every* method that requires `--id` as a flag would be able to be shortened to use a fifth argument which is implicitly interpreted as `id`.

