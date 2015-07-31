.. _services:

Services
========

Here you can find the services currently supported by Rackspace CLI, and those
services' supported commands.

Commands Reference
------------------

This section lists the commands you can use with Rackspace CLI.

.. toctree::
   :maxdepth: 3

   servers.rst
   files.rst
   blockstorage.rst
   networks.rst

Examples for services
---------------------

This section shows service-specific use cases for Rackspace CLI.

.. toctree::
   maxdepth: 2

   cloudfilesexamples.rst
   cloudserversexamples.rst
   cloudblockexamples.rst

Authenticating as another user
------------------------------

If you have an administrative account, you can authenticate as another user
with `--tenant-id` and `--auth-token` so that you can run commands with their
account information. This scenario is useful when you are troubleshooting for
someone.

Troubleshooting a request
-------------------------

If you want to see the request URL, response status or headers including
the X-Compute-Request-Id or X-Trans-Id, use the :code:`--log debug` option
on any `rack` command::

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
