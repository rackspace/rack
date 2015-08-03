.. _home:

Rackspace Command Line Interface
================================

.. warning:: The ``rack`` tool is under heavy development; the name of the binary
             and other portions may rapidly change. If you want to participate or
             provide early feedback, see the `Github project`_

Description
-----------

The Rackspace Command Line Interface is a unified tool to manage your Rackspace
services. It provides streamlined and secure configuration as well as a single
point of entry for all Rackspace Cloud services.


Quickstart
----------

For full instuctions on how to get started you should read :ref:`installation_and_configuration`.

The tl;dr version is to grab the binary for your platform:

* `Mac OSX (64 bit)`_
* `Linux (64 bit)`_
* `Windows (64 bit)`_

Once downloaded; you need to make it executable, if you are unfamiliar with this
or are running Windows, please see :ref:`installation_and_configuration`.

Next, run::

    rack configure

This command will automatically create a configuration file for you if it
doesn't exist and walk you through creating a profile for it::

    This interactive session will walk you through creating
    a profile in your configuration file. You may fill in all or none of the
    values.

    Rackspace Username: iamacat
    Rackspace API key: secrets
    Rackspace Region : IAD
    Profile Name (leave blank to create a default profile):

This allows you to immediately get working::

    rack servers instance list



Synopsis
--------

::

  rack <service> <subservice> <action> [--flags]

All ``rack`` commands follow the pattern above - for example, if you wanted to
list all running servers on your Rackspace account, you would type::

  rack servers instance list

And the response (**default**: table-based output) would look like::

      ID	Name		Status	Public IPv4	Private IPv4	Image	Flavor
      GUID	my_server	ACTIVE	101.130.19.31	10.208.128.233	GUID	io1-30


Options
-------

The ``rack`` CLI uses global options (``[--options]``) to alter the output, or
pass in other required **global** information into the tool, these are:

``--json``
  (boolean) Return output in JSON format.

``--table``
  Return output in tabular format. *This is the default output format.*

``--csv``
  Return output in csv format.

``--username``
  The Rackspace username to use for authentication.

``--api-key``
  The Rackspace API key to use for authentication.

``--region``
  The Rackspace region to use for authentication.

``--auth-url``
  The Rackspace URL to use for authentication. If not provided, this
  will default to the public U.S. Rackspace endpoint.

``--auth-tenant-id``
  The tenant ID to use for authentication. Prefixed with 'auth-' to avoid name
  collision with the ``tenant-id`` command flag.

``--auth-token``
  The token to use for authentication. This must be used with the ``auth-tenant-id``
  flag.

``--profile``
  The name of the config file profile to use to look for authentication credentials.

``--no-cache``
  Don't get or set authentication credentials in the rack cache.

``--no-header``
  Don't set the header for CSV nor tabular output.

``--use-service-net``
  Use the Rackspace internal URL to execute the request.

``--help, -h``
  Show help

For more on options, see :ref:`global_options`.

Services
--------

* :ref:`servers` - Commands for Rackspace Cloud Servers, dedicated and virtual.
* :ref:`files` - Commands for Rackspace Cloud Files.


.. toctree::
   :caption: Table of Contents
   :name: mastertoc
   :maxdepth: 2

   self
   configuration.rst
   globaloptions.rst
   services/index.rst



Indices and tables
==================

* :ref:`genindex`
* :ref:`search`

.. _Github project: https://github.com/jrperritt/rack
.. _Mac OSX (64 bit): https://ec4a542dbf90c03b9f75-b342aba65414ad802720b41e8159cf45.ssl.cf5.rackcdn.com/Darwin/amd64/rack
.. _Linux (64 bit): https://ec4a542dbf90c03b9f75-b342aba65414ad802720b41e8159cf45.ssl.cf5.rackcdn.com/Linux/amd64/rack
.. _Windows (64 bit): https://ec4a542dbf90c03b9f75-b342aba65414ad802720b41e8159cf45.ssl.cf5.rackcdn.com/Windows/amd64/rack.exe
.. _Cloud Control panel: https://mycloud.rackspace.com/
