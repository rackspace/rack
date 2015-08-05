.. _block_storage:

block-storage
=============

Commands for Rackspace Cloud Block Storage.

Synopsis
--------

::

   rack block-storage <resource> <action> [command flags]

Commands
--------

``volume``
~~~~~~~~~~

  Block Storage Volume operations

``list``
^^^^^^^^
Usage::

    rack block-storage volume list [optional flags]

``create``
^^^^^^^^^^
Usage::

    rack block-storage volume create --size <volumeSize> [optional flags]

``get``
^^^^^^^
Usage::

    rack block-storage volume get --id <volumeID> [optional flags]
    rack block-storage volume get --name <volumeName> [optional flags]
    (echo volumeID1 && echo volumeID2) | rack block-storage volume get --stdin id [optional flags]

``update``
^^^^^^^^^^
Usage::

    rack block-storage volume update --id <volumeID> [optional flags]
    rack block-storage volume update --name <volumeName> [optional flags]

``delete``
^^^^^^^^^^
Usage::

    rack block-storage volume delete --id <volumeID> [optional flags]
    rack block-storage volume delete --name <volumeName> [optional flags]
    (echo volumeID1 && echo volumeID2) | rack block-storage volume delete --stdin id [optional flags]

``snapshot``
~~~~~~~~~~~~

  Block Storage Snapshot operations

``list``
^^^^^^^^
Usage::

    rack block-storage snapshot list [optional flags]

``create``
^^^^^^^^^^
Usage::

    rack block-storage snapshot create --volume-id <volumeID> [optional flags]

``get``
^^^^^^^
Usage::

    rack block-storage snapshot get --id <snapshotID> [optional flags]
    rack block-storage snapshot get --name <snapshotName>] [optional flags]
    (echo snapshotID1 && echo snapshotID2) | rack block-storage snapshot get --stdin id [optional flags]

``delete``
^^^^^^^^^^
Usage::

    rack block-storage snapshot delete --id <snapshotID> [optional flags]
    rack block-storage snapshot delete --name <snapshotName> [optional flags]
    (echo snapshotID1 && echo snapshotID2) | rack block-storage snapshot delete --stdin id [optional flags]
