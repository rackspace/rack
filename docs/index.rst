.. _home:

Rackspace Command Line Interface
================================

The Rackspace Command Line Interface (``rack`` CLI) is a unified tool for managing your Rackspace
services. It provides streamlined and secure configuration and a single
point of entry for all Rackspace Cloud services.

Quickstart
----------

For full instructions, see :ref:`installation_and_configuration`.

To get started quickly, perform these steps:

1. Download the binary for your platform:

    * `Mac OS X (64-bit)`_
    * `Linux (64-bit)`_
    * `Windows (64-bit)`_

2. Make the binary executable. If you are unfamiliar with this process or you are running Windows, go to :ref:`installation_and_configuration`.

3. Go to the directory where you downloaded the ``rack`` binary and run the following command::

    ./rack configure

   This command automatically creates a configuration file if one doesn't already exist and walks you through creating a profile for it::

    This interactive session will walk you through creating
    a profile in your configuration file. You may fill in all or none of the
    values.

    Rackspace Username: <yourRackspaceUsername>
    Rackspace API key: <yourRackspaceApiKey>
    Rackspace Region: <yourRackspaceRegion>
    Profile Name (leave blank to create a default profile):

After you create the profile, you can immediately start working. For example, you could issue the following command to get a list of the servers on your Rackspace account::

    ./rack servers instance list

.. tip::
    Add the directory in which the ``rack`` binary resides to your system's PATH environment variable. Then, you can run it from anywhere, as shown in the following example::

      rack servers instance list


Overview
--------

All ``rack`` CLI commands use the following pattern::

    rack <command> <subcommand> <action> [flags]
    
.. note::
    ``<command>`` is usually equivalent to a Rackspace service.  

For example, to list all servers on your Rackspace account, you would issue the following command::

  rack servers instance list

The response, which is table-based output by default, would look as follows::

      ID	Name		Status	Public IPv4	Private IPv4	Image	Flavor
      GUID	my_server	ACTIVE	101.130.19.31	10.208.128.233	GUID	io1-30

Global Options
--------------

The ``rack`` CLI uses global options to alter output, authenticate, or
pass other global information into the tool. These options are
called *global* because they are valid for any command. To see a list of these options, go to :ref:`global_options`.

Services
--------

The ``rack`` CLI provides commands for interacting with the following Rackspace Cloud services:

* :ref:`servers` - Commands for Rackspace Cloud Servers, dedicated and virtual
* :ref:`files` - Commands for Rackspace Cloud Files
* :ref:`block_storage` - Commands for Rackspace Cloud Block Storage
* :ref:`networks` - Commands for Rackspace Cloud Networks
* :ref:`orchestration` - Commands for Rackspace Cloud Orchestration


.. _Github project: https://github.com/rackspace/rack
.. _Mac OS X (64-bit): https://ec4a542dbf90c03b9f75-b342aba65414ad802720b41e8159cf45.ssl.cf5.rackcdn.com/1.1.1/Darwin/amd64/rack
.. _Linux (64-bit): https://ec4a542dbf90c03b9f75-b342aba65414ad802720b41e8159cf45.ssl.cf5.rackcdn.com/1.1.1/Linux/amd64/rack
.. _Windows (64-bit): https://ec4a542dbf90c03b9f75-b342aba65414ad802720b41e8159cf45.ssl.cf5.rackcdn.com/1.1.1/Windows/amd64/rack.exe
.. _Cloud Control panel: https://mycloud.rackspace.com/
