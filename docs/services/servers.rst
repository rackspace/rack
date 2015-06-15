.. _servers:

servers
=======

Commands for Rackspace Cloud Servers, dedicated and virtual.

Synopsis
--------

::

   rack servers command [command options] [arguments...]

Commands
--------

``instance``
~~~~~~~~~~~~

  Server Instance operations

``list``
^^^^^^^^
Example Usage::

    rack servers instance list [optional flags]

``create``
^^^^^^^^^^
Example Usage::

    rack servers instance create [--name <serverName>] [optional flags]

``get``
^^^^^^^
Example Usage::

    rack servers instance get [--id <serverID> | --name <serverName>] [optional flags]

``update``
^^^^^^^^^^
Example Usage::

    rack servers instance update [--id <serverID> | --name <serverName>] [optional flags]

``delete``
^^^^^^^^^^
Example Usage::

    rack servers instance delete [--id <serverID> | --name <serverName>] [optional flags]

``reboot``
^^^^^^^^^^
Example Usage::

    rack servers instance reboot [--id <serverID> | --name <serverName>] [--soft | --hard] [optional flags]

``rebuild``
^^^^^^^^^^^
Example Usage::

    rack servers instance rebuild [--id <serverID> | --name <serverName>] [--imageID <imageID>] [--adminPass <adminPass>] [optional flags]

``resize``
^^^^^^^^^^
Example Usage::

    rack servers instance resize [--id <serverID> | --name <serverName>] [--flavorID <flavorID>] [optional flags]


``image``
~~~~~~~~~

  Server Image operations

``list``
^^^^^^^^
Example Usage::

    rack servers image list [flags]

``get``
^^^^^^^^
Example Usage::

    rack servers image get <imageID> [flags]

``flavor``
~~~~~~~~~~

  Server Flavor operations

``list``
^^^^^^^^
Example Usage::

    rack servers flavor list [flags]

``get``
^^^^^^^
Example Usage::

    rack servers flavor get <flavorID> [flags]


``keypair``
~~~~~~~~~~~

  Server Keypair operations

``list``
^^^^^^^^
Example Usage::

    rack servers keypair list [flags]

``create``
^^^^^^^^^^
Example Usage::

    rack servers keypair create <keypairName> [flags]

``get``
^^^^^^^
Example Usage::

    rack [globals] servers keypair get [--name <keypairName>] [flags]

``delete``
^^^^^^^^^^
Example Usage::

    rack servers keypair delete [--name <keypairName>] [flags]
