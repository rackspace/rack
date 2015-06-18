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
Usage::

    rack servers instance list [optional flags]

``create``
^^^^^^^^^^
Usage::

    rack servers instance create [--name <serverName>] [optional flags]

``get``
^^^^^^^
Usage::

    rack servers instance get [--id <serverID> | --name <serverName>] [optional flags]

``update``
^^^^^^^^^^
Usage::

    rack servers instance update [--id <serverID> | --name <serverName>] [optional flags]

``delete``
^^^^^^^^^^
Usage::

    rack servers instance delete [--id <serverID> | --name <serverName>] [optional flags]

``reboot``
^^^^^^^^^^
Usage::

    rack servers instance reboot [--id <serverID> | --name <serverName>] [--soft | --hard] [optional flags]

``rebuild``
^^^^^^^^^^^
Usage::

    rack servers instance rebuild [--id <serverID> | --name <serverName>] --imageID <imageID> --adminPass <adminPass> [optional flags]

``resize``
^^^^^^^^^^
Usage::

    rack servers instance resize [--id <serverID> | --name <serverName>] --flavorID <flavorID> [optional flags]


``image``
~~~~~~~~~

  Server Image operations

``list``
^^^^^^^^
Usage::

    rack servers image list [flags]

``get``
^^^^^^^^
Usage::

    rack servers image get [--id <serverID> | --name <serverName>] [optional flags]

``flavor``
~~~~~~~~~~

  Server Flavor operations

``list``
^^^^^^^^
Usage::

    rack servers flavor list [flags]

``get``
^^^^^^^
Usage::

    rack servers flavor get [--id <serverID> | --name <serverName>] [optional flags]


``keypair``
~~~~~~~~~~~

  Server Keypair operations

``list``
^^^^^^^^
Usage::

    rack servers keypair list [flags]

``create``
^^^^^^^^^^
Usage::

    rack servers keypair create --name <keypairName> [optional flags]

``get``
^^^^^^^
Usage::

    rack [globals] servers keypair get --name <keypairName> [optional flags]

``delete``
^^^^^^^^^^
Usage::

    rack servers keypair delete --name <keypairName> [optional flags]
