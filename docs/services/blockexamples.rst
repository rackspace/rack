.. _blockexamples

============================
Cloud Block Storage examples
============================

Before you get started, be sure you have entered your username and API key
and stored them locally::

    rack configure

You can get help for any commands and its options by adding `--help` to the
series of commands::

    $ rack block-storage snapshot create --help

Create a volume
---------------

Create a volume by using the `rack block-storage volume create --size`
command, than specifying the size and name you wish the volume to be. In the
example, a volume named "Store" is created::

    $ rack block-storage volume create --size 75 --name Store
    ID		81c7a7e5-01a5-44bb-9b43-0cc9f7c4e423
    Name
    Description
    Size		75
    VolumeType	SATA
    SnapshotID
    Attachments	[]
    CreatedAt	2015-07-31T17:49:41.752349

You can also attach the new volume to a server using the `volume-id`. In the example,
the volume above is attached to server named "RACK"::

    $ rack servers volume-attachment create --server-name RACK --volume-id 81c7a7e5-01a5-44bb-9b43-0cc9f7c4e423
    ID	81c7a7e5-01a5-44bb-9b43-0cc9f7c4e423
    Device	/dev/vvdc
    VolumeID81c7a7e5-01a5-44bb-9b43-0cc9f7c4e423
    ServerID8a254ea3-77b5-4f74-a893-8d2d51ae2cca

Create a snapshot
-----------------

Create a snapshot by detaching detaching the volume from it's server. You will
need to use the volume's `attachment id`, which can be found with the
`rack server volume-attachment list` command::

    $ rack servers volume-attachment list --server-name RACK
    ID					                          Device		VolumeID				                      ServerID
    81c7a7e5-01a5-44bb-9b43-0cc9f7c4e423	/dev/xvdc	81c7a7e5-01a5-44bb-9b43-0cc9f7c4e423	8a254ea3-77b5-4f74-a893-8d2d51ae2cca

Use the command `rack servers volume-attach delete` followed by the `attachment id` and
the server name::

    $ rack servers volume-attach delete --server-name RACK --id 81c7a7e5-01a5-44bb-9b43-0cc9f7c4e423

Once the server is detached, you can now create a snapshot using the command
`rack block-storage snapshot create`, followed by `volume-id` and the name
you wish the snapshot be. In the example, a snapshot called "Store1" is created::

    $ rack block-storage snapshot create --volume-id 81c7a7e5-01a5-44bb-9b43-0cc9f7c4e423
    ID		180a6c5c-ad6d-4cb6-846f-d500d67e59a5
    Name	Store1
    Description
    Size		75
    VolumeType
    SnapshotID
    Attachments
    CreatedAt	2015-07-31T18:57:34.652136
