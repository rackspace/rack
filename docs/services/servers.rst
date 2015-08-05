.. _servers:

=======
Servers
=======

This section lists all cloud servers commands supported by Rackspace CLI.

Commands
--------

All cloud server commands are based on this syntax::

   rack servers <resource> <action> [command flags]

``instance``
~~~~~~~~~~~~

Cloud server instance commands use this syntax::

    rack servers instance <action> [command flags]

``list``
^^^^^^^^
Retrieves a list of virtual and bare metal servers::

    rack servers instance list [optional flags]

``create``
^^^^^^^^^^
Creates a server instance::

    rack servers instance create --name <serverName> [optional flags]
    (echo serverName1 && echo serverName2) | rack servers instance create --stdin name [optional flags]

``get``
^^^^^^^
Retrieves details of a specified server::

    rack servers instance get --id <serverID> [optional flags]
    rack servers instance get --name <serverName> [optional flags]
    (echo serverID1 && echo serverID2) | rack servers instance get --stdin id [optional flags]

``update``
^^^^^^^^^^
Updates one or more editable attributes of a specified server
instance::

    rack servers instance update --id <serverID> [optional flags]
    rack servers instance update --name <serverName> [optional flags]

``delete``
^^^^^^^^^^
Deletes a server instance::

    rack servers instance delete --id <serverID> [optional flags]
    rack servers instance delete --name <serverName> [optional flags]
    (echo serverID1 && echo serverID2) | rack servers instance delete --stdin id [optional flags]

``reboot``
^^^^^^^^^^
Performs a soft or hard reboot of a specified server. A soft reboot
will slowly shutdown and restart your server's operating system. A hard reboot
will perform an immediate shutdown and restart::

    rack servers instance reboot --id <serverID> [--soft | --hard] [optional flags]
    rack servers instance reboot --name <serverName> [--soft | --hard] [optional flags]
    (echo serverID1 && echo serverID2) | rack servers instance reboot --stdin id [--soft | --hard] [optional flags]

``rebuild``
^^^^^^^^^^^
Removes all data on the server and replaces it with the specified image::

    rack servers instance rebuild --id <serverID> --image-id <imageID> --admin-pass <adminPass> [optional flags]
    rack servers instance rebuild --name <serverName> --image-id <imageID> --admin-pass <adminPass> [optional flags]

``resize``
^^^^^^^^^^
Converts an existing server to a different flavor, which scales the
server up or down. The original server is saved for a period of time to allow roll
back if a problem occurs::

    rack servers instance resize --id <serverID> --flavor-id <flavorID> [optional flags]
    rack servers instance resize --name <serverName> --flavor-id <flavorID> [optional flags]
    (echo serverID1 && echo serverID2) | rack servers instance resize --stdin id --flavor-id <flavorID> [optional flags]

.. note::
    This command is not available for OnMetal servers.

``set-metadata``
^^^^^^^^^^^^^^^^
Sets metadata for the specified server or image::

    rack servers instance set-metadata --id <serverID> --metadata <key1=val1,key2=val2,...> [optional flags]
    rack servers instance set-metadata --name <serverName> --metadata <key1=val1,key2=val2,...> [optional flags]

``get-metadata``
^^^^^^^^^^^^^^^^
Retrieves the metadata for a given server::

    rack servers instance get-metadata --id <serverID> [optional flags]
    rack servers instance get-metadata --name <serverName> [optional flags]

``update-metadata``
^^^^^^^^^^^^^^^^^^^
Updates metadata items for a specified server or image, or adds the specified
metadata if there is no current metadata associated with the server or image::

    rack servers instance update-metadata --id <serverID> --metadata <key1=val1,key2=val2,...> [optional flags]
    rack servers instance update-metadata --name <serverName> --metadata <key1=val1,key2=val2,...> [optional flags]

``delete-metadata``
^^^^^^^^^^^^^^^^^^^
Delete one or more metadata keys from a server::

    rack servers instance delete-metadata --id <serverID> --metadata-keys <key1,key2,...> [optional flags]
    rack servers instance delete-metadata --name <serverName> --metadata-keys <key1,key2,...> [optional flags]


``image``
~~~~~~~~~

Cloud server image commands use this syntax::

    rack server image <action> [optional flags]

``list``
^^^^^^^^
Lists all images visible by your account::

    rack servers image list [optional flags]

``get``
^^^^^^^
Retrieves details of the specified image::

    rack servers image get --id <imageID> [optional flags]
    rack servers image get --name <imageName>] [optional flags]
    (echo imageID1 && echo imageID2) | rack servers image get --stdin id [optional flags]

.. note::

   To guarantee usage of the same image every time, use the `id` flag. Images often
   are updated with security patches, and the updated images will have a different ID but
   the same name.


``flavor``
~~~~~~~~~~

Cloud server flavor commands use this syntax::

    rack servers flavor <action> [optional flags]

``list``
^^^^^^^^
Lists information for all available flavors::

    rack servers flavor list [optional flags]

``get``
^^^^^^^
Retrieves details of the specified flavor::

    rack servers flavor get --id <flavorID> [optional flags]
    rack servers flavor get --name <flavorName>] [optional flags]
    (echo flavorID1 && echo flavorID2) | rack servers flavor get --stdin id [optional flags]

``keypair``
~~~~~~~~~~~

Server keypair commands use this syntax::

    rack servers keypair <action> [optional flags]

``list``
^^^^^^^^
Retrieves a list of all key pairs associated with your account::

    rack servers keypair list [flags]

``generate``
^^^^^^^^^^^^
Generates a newly created key pair with the specified name::

    rack servers keypair generate --name <keypairName> [optional flags]
    (echo keypairName1 && echo keypairName2) | rack servers keypair generate --stdin name [optional flags]

``upload``
^^^^^^^^^^
Uploads an existing key pair with the specified name::

    rack servers keypair upload --name <keypairName> --public-key <publicKeyData> [optional flags]
    rack servers keypair upload --name <keypairName> --file <publicKeyfile> [optional flags]

``get``
^^^^^^^
Retrieves details on a specified key pair::

    rack [globals] servers keypair get --name <keypairName> [optional flags]
    (echo keypairName1 && echo keypairName2) | rack servers keypair get --stdin name [optional flags]

``delete``
^^^^^^^^^^
Deletes the specified key paid::

    rack servers keypair delete --name <keypairName> [optional flags]
    (echo keypairName1 && echo keypairName2) | rack servers keypair delete --stdin name [optional flags]


``volume-attachment``
~~~~~~~~~~~~~~~~~~~~~

Cloud server volume attachment commands use this syntax::

    rack server volume-attachment <action> [optional flag]

These commands are often used with :ref:`cloud block storage <blockexamples>`.

``list``
^^^^^^^^
Lists the volume attachments for the specified server::

    rack servers volume-attachment list --server-id <serverID> [optional flags]
    rack servers volume-attachment list --server-name <serverName> [optional flags]
    rack servers volume-attachment list --stdin server-id [optional flags]

``create``
^^^^^^^^^^
Attaches one or more volumes to the specified sever::

    rack servers volume-attachment create --server-id <serverID> --volume-id <volumeID> [optional flags]
    rack servers volume-attachment create --server-name <serverName> --volume-id <volumeID> [optional flags]
    rack servers volume-attachment create --server-id <serverID> --volume-name <volumeName> [optional flags]
    rack servers volume-attachment create --server-name <serverName> --volume-name <volumeName> [optional flags]
    (echo volumeID1 && echo volumeID2) | rack servers volume-attachment create --server-id <serverID> --stdin volume-id [optional flags]
    (echo volumeID1 && echo volumeID2) | rack servers volume-attachment create --server-name <serverName> --stdin volume-id [optional flags]

``get``
^^^^^^^
Retrieves details of a specified volume attachment ID for a specified server::

    rack servers volume-attachment get --server-id <serverID> --id <attachmentID> [optional flags]
    rack servers volume-attachment get --server-name <serverName> --id <attachmentID> [optional flags]

``delete``
^^^^^^^^^^
Deletes a specified volume attachment from a specified server instance::

    rack servers volume-attachment delete --server-id <serverID> --id <attachmentID> [optional flags]
    rack servers volume-attachment delete --server-name <serverName> --id <attachmentID> [optional flags]
