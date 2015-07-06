.. _global_options:

Global Options
==============

Global options are command line flags that precede any service call or command
and allow you to customize various aspects of the ``rack`` tool at runtime.
These options may override configuration file or environment variables you have
set previously, change output format or other aspects of the tool.

For example:

::

    rack --json <service> <command> <subcommand> [--flags]

Would result in all commands returning a JSON_ formatted output.

Options
-------

``--json``
~~~~~~~~~~

  (boolean) Return output in JSON format.

When added to the arguments; output changes from the default table-based output
to JSON.

Given::


    rack servers instance list
    ID	        Name		Status	Public IPv4	Private IPv4	Image	Flavor
    GUID	my_server	ACTIVE	101.130.19.31	10.208.128.233	GUID	io1-30

Adding the JSON option returns::

    rack --json servers instance list
    [
      {
        "Flavor": "io1-30",
        "ID": "GUID",
        "Image": "fd1b9e98-21c5-43bb-89d7-bfc4ee3c5caf",
        "Name": "my_server",
        "Private IPv4": "10.208.128.233",
        "Public IPv4": "101.130.19.31",
        "Status": "ACTIVE"
      }
    ]

When the output pipe is **not** a tty; the JSON is no longer "pretty printed" and
can be used when passing straight into other commands that require a JSON_
payload or argument or another service.

``--table``
~~~~~~~~~~~
  (boolean) Return output in tabular format.

Default output format for ``rack``.

Given::

    rack servers instance list
    ID	        Name		Status	Public IPv4	Private IPv4	Image	Flavor
    GUID	my_server	ACTIVE	101.130.19.31	10.208.128.233	GUID	io1-30

This presents a well formatted table with headers.

You can add the ``--table`` option if you have set defaults to JSON, CSV, etc
elsewhere.

``--csv``
~~~~~~~~~

  (boolean) Return output in csv format.

CSV, or comma separated output is useful for passing to other operating system
tools, importing into Excel, Google Sheets, or another data tool.

Given::

    rack servers instance list
    ID	        Name		Status	Public IPv4	Private IPv4	Image	Flavor
    GUID	my_server	ACTIVE	101.130.19.31	10.208.128.233	GUID	io1-30

Adding the CSV option returns::

    rack --csv servers instance list
    ID,Name,Status,Public IPv4,Private IPv4,Image,Flavor
    GUID,my_server,ACTIVE,101.130.19.31,10.208.128.233,fd1b9e98-21c5-43bb-89d7-bfc4ee3c5caf,io1-30

This presents a compact format with appropriate CSV headers.

``--log``
~~~~~~~~~

  (string) Log relevant information about the HTTP request. Options are: info, debug.

  Example: ``rack servers keypair list --log info``

``--username``
~~~~~~~~~~~~~~

  (string) The Rackspace username to use for authentication.

``--apikey``
~~~~~~~~~~~~

  (string) The Rackspace API key to use for authentication.

``--region``
~~~~~~~~~~~~

  (string) The Rackspace region to use for authentication.

``--authurl``
~~~~~~~~~~~~~

  (string) The Rackspace URL to use for authentication. If not provided, this
  will default to the public U.S. Rackspace endpoint.

``--profile``
~~~~~~~~~~~~~

  (string) The name of the config file profile to use to look for authentication credentials.

``--no-cache``
~~~~~~~~~~~~~~

  (boolean) Don't get or set authentication credentials in the rack cache.

``--help, -h``
~~~~~~~~~~~~~~

  (boolean) Show help in a given context.

Help is available on the base level; for example::

    rack help
    NAME:
       rack - An opinionated CLI for the Rackspace cloud

    USAGE:
       rack [global options] command [command options] [arguments...]

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
       --generate-bash-completion
       --version, -v		print the version

And it is available per service::

    rack servers help
    NAME:
       rack servers - Used for the Servers service

    USAGE:
       rack servers [global options] command [command options] [arguments...]

    VERSION:
       0.0.0

    COMMANDS:
       instance	Used for Server Instance operations
       image	Used for Server Image operations
       flavor	Used for Server Flavor operations
       keypair	Used for Server Keypair operations
       help, h	Shows a list of commands or help for one command


And again, per command:

    rack servers keypair help
    NAME:
       rack servers keypair - Used for Server Keypair operations

    USAGE:
       rack servers keypair [global options] command [command options] [arguments...]

    VERSION:
       0.0.0

    COMMANDS:
       list		rack servers keypair list [flags]
       create	rack servers keypair create <keypairName> [flags]
       get		rack [globals] servers keypair get [--name <keypairName>] [flags]
       delete	rack servers keypair delete [--name <keypairName>] [flags]
       help, h	Shows a list of commands or help for one command

``--version, -v``
~~~~~~~~~~~~~~~~~

  Print the version of the ``rack`` CLI.

The version number of the CLI will be important when opening tickets, filing
issues on the issue tracker or in any other debugging session. Please include
this any time you are having issues.


.. JSON: http://json.org/
