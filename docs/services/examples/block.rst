.. _blockexamples

======================
Block Storage examples
======================

Before you start using examples, be sure to enter your username and API key and store them locally by running the ``rack configure`` command. For more information, see :ref:`installation_and_configuration`.

You can get help for any command and its options by adding ``--help`` to the series of commands::

    $ rack block-storage snapshot create --help

Create a volume
~~~~~~~~~~~~~~~

Create a volume by using the ``rack block-storage volume create --size`` command and specifying the size and name for the volume. In the following example, a volume named ``Store`` is created.

::

    $ rack block-storage volume create --size 75 --name Store
    ID		81c7a7e5-01a5-44bb-9b43-0cc9f7c4e423
    Name
    Description
    Size		75
    VolumeType	SATA
    SnapshotID
    Attachments	[]
    CreatedAt	2015-07-31T17:49:41.752349

You can also attach the new volume to a server by using the ``volume-id`` flag. In the following example, the newly created volume is attached to a server named ``RACK``.

::

    $ rack servers volume-attachment create --server-name RACK --volume-id 81c7a7e5-01a5-44bb-9b43-0cc9f7c4e423
    ID	81c7a7e5-01a5-44bb-9b43-0cc9f7c4e423
    Device	/dev/vvdc
    VolumeID81c7a7e5-01a5-44bb-9b43-0cc9f7c4e423
    ServerID8a254ea3-77b5-4f74-a893-8d2d51ae2cca

Create a snapshot
~~~~~~~~~~~~~~~~~

Create a snapshot by detaching the volume from its server. You use the volume's attachment ID, which you can find by using the following command.

::

    $ rack servers volume-attachment list --server-name RACK
    ID					                          Device		VolumeID				                      ServerID
    81c7a7e5-01a5-44bb-9b43-0cc9f7c4e423	/dev/xvdc	81c7a7e5-01a5-44bb-9b43-0cc9f7c4e423	8a254ea3-77b5-4f74-a893-8d2d51ae2cca

Use the following command, followed by the server name and the attachment ID, to detach the volume from the server. 

::

    $ rack servers volume-attachment delete --server-name RACK --id 81c7a7e5-01a5-44bb-9b43-0cc9f7c4e423

After the volume is detached, you can create a snapshot by using the following command and specifying the ID of the volume and a name for the snapshot. In the following example, a snapshot called ``Store1`` is created.

::

    $ rack block-storage snapshot create --volume-id 81c7a7e5-01a5-44bb-9b43-0cc9f7c4e423
    ID		180a6c5c-ad6d-4cb6-846f-d500d67e59a5
    Name	Store1
    Description
    Size		75
    VolumeType
    SnapshotID
    Attachments
    CreatedAt	2015-07-31T18:57:34.652136
