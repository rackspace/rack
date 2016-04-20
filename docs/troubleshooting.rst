.. _troubleshooting:

===============
Troubleshooting
===============

This section provides suggestions for troubleshooting issues when you use the Rackspace CLI.

Troubleshooting a request
-------------------------

If you want to see a request URL, response status, or headers, including the ``X-Compute-Request-Id`` or ``X-Trans-Id`` header, use the ``--log debug`` option on any ``rack`` command::

    $ rack servers instance list --fields publicipv4 --log debug

    INFO[0000] Global Options:
    output: table (from default value)
    log: debug (from command-line)


    INFO[0000] Authentication Credentials:
    region: ORD (from config file (profile: DEFAULT))
    auth-url: https://identity.api.rackspacecloud.com/v2.0/ (from default value)
    username: agent (from config file (profile: DEFAULT))
    api-key: ccbbccbbccbbccbbccbbccbbccbbaaac (from config file (profile: DEFAULT))


    INFO[0000] Using token from cache: 9c99727c0af4434ea41627761683a012

    INFO[0000] Request URL: https://ord.servers.api.rackspacecloud.com/v2/1234567/servers/detail

    DEBU[0000] Response Status: 200 OK

    DEBU[0000] Response Headers: {
      "Content-Length": [
        "3087"
      ],
      "Content-Type": [
        "application/json"
      ],
      "Date": [
        "Thu, 30 Jul 2015 19:58:23 GMT",
        "Thu, 30 Jul 2015 19:58:23 GMT"
      ],
      "Server": [
        "Jetty(9.2.z-SNAPSHOT)"
      ],
      "Via": [
        "1.1 Repose (Repose/7.1.3.0)"
      ],
      "X-Compute-Request-Id": [
        "req-edb79c8c-9eaf-4314-b055-59276172292c"
      ],
      "X-Trans-Id": [
        "6a52a395-4f59-4312-a852-0826d11b4a0c"
      ]
    }

    PublicIPv4
    162.209.0.32
    23.253.213.33

.. _authenticating:

Authenticating as another user
------------------------------

If you have an administrative account, you can authenticate as another user with the ``--tenant-id`` and ``--auth-token`` flags to run commands with that user's account information. This scenario is useful when you are troubleshooting for someone else.

Error when piping results to another command
--------------------------------------------

When you want to use the results of one ``rack`` command to drive input for another ``rack`` command, be sure to remove the header from the table listing output by using the ``--no-header`` option. Otherwise you might see a series of errors because the header is considered a request. 

Following is an example of what *not* to do::

    $ rack files object list --container temp --fields name | rack files object delete --container temp --stdin name

The result is the following error::

    error    I couldn't find object [Name] in container [temp]

Instead, run the command with the ``--no-header`` option::

    $ rack files object list --container temp --fields name --no-header | rack files object delete --container temp --stdin name
