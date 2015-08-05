.. _block_storage:

<<<<<<< HEAD
=============
Block storage
=======
block-storage
>>>>>>> fb9b3c09d96ce50bb2117a52db8c7d591d6d5157
=============

This section lists all block storage commands supported by Rackspace CLI

Commands
--------

All block storage commands are based on this syntax::

   rack block-storage <resource> <action> [command flags]

``volume``
~~~~~~~~~~

Block storage volume commands use this syntax::

    rack block-storage volume <action> [optional flags]

``list``
^^^^^^^^
Retrieves a list of volumes::

    rack block-storage volume list [optional flags]

``create``
^^^^^^^^^^
Creates a volume::

    rack block-storage volume create --size <volumeSize> [optional flags]

``get``
^^^^^^^
Retrieves details on a specified volume::

    rack block-storage volume get --id <volumeID> [optional flags]
    rack block-storage volume get --name <volumeName> [optional flags]
    (echo volumeID1 && echo volumeID2) | rack block-storage volume get --stdin id [optional flags]

``update``
^^^^^^^^^^
Updates the name and description of a volume::

    rack block-storage volume update --id <volumeID> [optional flags]
    rack block-storage volume update --name <volumeName> [optional flags]

``delete``
^^^^^^^^^^
Permanently removes a volume::

    rack block-storage volume delete --id <volumeID> [optional flags]
    rack block-storage volume delete --name <volumeName> [optional flags]
    (echo volumeID1 && echo volumeID2) | rack block-storage volume delete --stdin id [optional flags]

``snapshot``
~~~~~~~~~~~~

Block storage snapshot commands use this syntax::

    rack block-storage snapshot <actions> [optional flags]

``list``
^^^^^^^^
Retrieves a list of snapshots::

    rack block-storage snapshot list [optional flags]

``create``
^^^^^^^^^^
Creates a snapshot based on a specified volume id::

    rack block-storage snapshot create --volume-id <volumeID> [optional flags]

``get``
^^^^^^^
Retrieves details on a specified snapshot::

    rack block-storage snapshot get --id <snapshotID> [optional flags]
    rack block-storage snapshot get --name <snapshotName>] [optional flags]
    (echo snapshotID1 && echo snapshotID2) | rack block-storage snapshot get --stdin id [optional flags]

``delete``
^^^^^^^^^^
Permanently removes a snapshot::

    rack block-storage snapshot delete --id <snapshotID> [optional flags]
    rack block-storage snapshot delete --name <snapshotName> [optional flags]
    (echo snapshotID1 && echo snapshotID2) | rack block-storage snapshot delete --stdin id [optional flags]
