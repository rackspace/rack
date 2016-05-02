.. _global_options:

Global options
==============

Global options are command-line flags that are valid for any command and enable you to customize various aspects of the ``rack`` CLI tool at runtime. These options might override configuration file or environment variables that you have set previously, or change the output format or other aspects of the tool.

For example, you can get JSON_ formatted output from any series of commands by appending the ``--output json`` flag::

    rack <command> <subcommand> <action> --output json
    
You can display any available command-line flags by appending ``--help`` to the command::

    rack <command> <subcommand> <action> --help

Options
-------

This section describes the global options available in the ``rack`` CLI. 

``--output``
~~~~~~~~~~~~

String. The format in which to return the output. Options are ``json``, ``table``, and ``csv``. Default is ``table``.

``json``
^^^^^^^^

Returns output in JSON.

::

    rack servers instance list --output json
    [
      {
        "Flavor": "io1-30",
        "ID": "GUID",
        "Image": "GUID",
        "Name": "my_server",
        "Private IPv4": "10.208.128.233",
        "Public IPv4": "101.130.19.31",
        "Status": "ACTIVE"
      }
    ]

.. note::
    When the output pipe is *not* a terminal, the JSON is not "pretty printed" and can be used when passing straight into other commands that require a JSON_ payload.

``table``
^^^^^^^^^

Returns output in a well-formatted table with headers. This is the default output format for ``rack``.

::

    rack servers instance list
    ID	        Name		Status	Public IPv4	Private IPv4	Image	Flavor
    GUID	my_server	ACTIVE	101.130.19.31	10.208.128.233	GUID	io1-30

You can add the ``--output table`` option if you have set defaults to ``json`` or ``csv`` elsewhere. You can use the ``--no-header`` option to return a table without headers.

``csv``
^^^^^^^

Returns output in comma-separated value (CSV) format. The format is compact with appropriate CSV headers.

The CSV format is useful for passing data to other operating system tools and importing data into Excel, Google Sheets, or another data tool.

::

  rack servers instance list --output csv
  ID,Name,Status,Public IPv4,Private IPv4,Image,Flavor
  GUID,my_server,ACTIVE,101.130.19.31,10.208.128.233,GUID,io1-30

``--log``
~~~~~~~~~

String. Displays log information at the info or debug level for the command. Options are ``info`` and ``debug``.

::

    rack servers keypair list --log info

``--username``
~~~~~~~~~~~~~~

String. The Rackspace username to use for authentication.

``--api-key``
~~~~~~~~~~~~~

String. The Rackspace API key to use for authentication.

``--auth-tenant-id``
~~~~~~~~~~~~~~~~~~~~

String. The tenant ID to use for authentication. This option can be provided only as a command-line flag. It is prefixed with ``auth-`` so that it does not collide with the ``tenant-id`` command flags.

``--auth-token``
~~~~~~~~~~~~~~~~

String. The token to use for authentication. This option can be provided only as a command-line flag. It must be used with the ``auth-tenant-id`` flag.

``--region``
~~~~~~~~~~~~

String. The Rackspace region to use for authentication.

``--auth-url``
~~~~~~~~~~~~~~

String. The Rackspace URL to use for authentication. If this option is not provided, the URL defaults to the public U.S. Rackspace endpoint.

``--profile``
~~~~~~~~~~~~~

String. The name of the profile (in the configuration file) to use to look for authentication credentials.

``--no-cache``
~~~~~~~~~~~~~~

Boolean. Indicates not to get or set authentication credentials in the ``rack`` cache.

``--no-header``
~~~~~~~~~~~~~~~

Boolean. Indicates not to set the header for CSV or tabular output. This option is helpful when you are piping output from a ``list`` command.

``--use-service-net``
~~~~~~~~~~~~~~~~~~~~~

Boolean. Indicates to use the Rackspace internal URL to execute the request. This option is useful only when you are running a ``rack`` command from a Rackspace server.

``--help, -h``
~~~~~~~~~~~~~~

Boolean. Shows help in a given context.

Example of help at the base level::

    rack --help
    NAME:
       rack - An opinionated CLI for the Rackspace cloud

    USAGE:
       rack <service> <subservice> <action> [flags]

    VERSION:
       0.0.0

    COMMANDS:
       servers	Used for the Servers service
       help, h	Shows a list of commands or help for one command

    GLOBAL OPTIONS:
       --json			Return output in JSON format.
       --table			Return output in tabular format. This is the default output format.
       --csv			Return output in csv format.
       --help, -h			show help

Example of help at the command (service) level::

    rack servers --help
    NAME:
       rack servers - Used for the Servers service

    USAGE:
       rack servers <subservice> <action> [flags]

    VERSION:
       0.0.0

    COMMANDS:
       instance	Used for Server Instance operations
       image	Used for Server Image operations
       flavor	Used for Server Flavor operations
       keypair	Used for Server Keypair operations
       help, h	Shows a list of commands or help for one command


Example of help at the subcommand (subservice) level::

    rack servers keypair --help
    NAME:
       rack servers keypair - Used for Server Keypair operations

    USAGE:
       rack servers keypair <action> [flags]

    VERSION:
       0.0.0

    COMMANDS:
       list		rack servers keypair list [flags]
       create	rack servers keypair create <keypairName> [flags]
       get		rack [globals] servers keypair get [--name <keypairName>] [flags]
       delete	rack servers keypair delete [--name <keypairName>] [flags]
       help, h	Shows a list of commands or help for one command

Example of help at the action level::

    rack servers instance list --help
    NAME: list - rack servers instance list  [flags]

    DESCRIPTION: Lists existing servers

    COMMAND FLAGS:
    --all-pages     [optional] Return all servers. Default is to paginate.
    --name          [optional] Only list servers with this name.
    --changes-since [optional] Only list servers that have been changed since this time/date stamp.
    --image         [optional] Only list servers that have this image ID.
    --flavor        [optional] Only list servers that have this flavor ID.
    --status        [optional] Only list servers that have this status.
    --marker        [optional] Start listing servers at this server ID.
    --limit         [optional] Only return this many servers at most.
    --fields        [optional] Only return these comma-separated case-insensitive fields.
                    Choices: id, name, status, publicipv4, privateipv4, image, flavor


    GLOBAL FLAGS:
    --username              The username with which to authenticate.
    --api-key               The API key with which to authenticate.
    --auth-tenant-id        The tenant ID of the user to authenticate as. May only be provided as a command-line flag.
    --auth-token            The authentication token of the user to authenticate as. This must be used with the `auth-tenant-id` flag.
    --auth-url              The endpoint to which authenticate.
    --region                The region to which authenticate.
    --use-service-net       Whether or not to use the internal Rackspace network
    --profile               The config file profile to use for authentication.
    --output                Format in which to return output. Options: json, csv, table. Default is 'table'.
    --no-cache              Don't get or set authentication credentials in the rack cache.
    --log                   Print debug information from the command. Options are: debug, info
    --no-header             Don't return a header for CSV nor tabular output.


.. JSON: http://json.org/
