# Rack CLI Design Document

This document will attempt to outline the design for the Rackspace Global Command Line Interface. It will not exhaustively define each interface; rather, it will focus on the philosophy and representative examples that should be used to inform the successive specifications for the specific service implementations.

As such, there are 4 core tenants we're using to define the design.

1. Correctness
2. Consistency
3. Completeness
4. Terseness

## Correctness

As a top level objective, the CLI must operate correctly. It's unacceptable to deliver a non-functional tool, nor one that behaves contrary to the documented expectations. The key here is to document the expected behavior, and then ensure that the delivered tool matches those expectations.

## Consistency

Our ambition for consistency is to provide a model that is used throughout the command line interface. Ideally, all calls take format:

```
rack <service> [--options] <model> <action> [--flags]
```

An example of `service` might be `servers` or `files` or `networks`. `model` might be `instance`, `object` or `container`, or perhaps `network`. Action would be the underlying API call to make, for example `create` or `delete`.

All of the following are examples of the pattern being applied consistently:

* `rack servers instances list` - Get a list of all instances
* `rack files containers list` - Get a list of all containers
* `rack networks security-groups get --id 12345` - Get a security group by id
* `rack load-balancers node remove --id 12345 --load-balancer-id abcdef` - Remove a node from a load balancer

## Completeness

While versions of `rack` may be shipped without complete coverage of all services, our goal is to have 100% coverage of services and capabilities within those services of the Rackspace public cloud. Additionally, services should not be released without complete functional support within that service. For example, shipping Cloud Files support without container operations would not fulfill our expectations.

## Terseness

To the extent of the preceding objectives, terseness is the last priority. We should, so long as it does not come at the expense of the goals listed above, optimize the command line interface for efficiency.

For example, we may be able to imply values without the `--id` flag as appropriate:

```
rack servers instances get --id 12345
```

might be synonymous with:

```
rack servers instances get 12345
```

This is an example of ways to add terseness without compromising consistency. The key in this example would be that *every* method that requires `--id` as a flag would be able to be shortened to use a fifth argument which is implicitly interpreted as `id`.

