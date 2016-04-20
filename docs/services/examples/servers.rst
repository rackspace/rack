.. _serversexamples:

================
Servers examples
================

Before you start using examples, be sure to enter your username and API key and store them locally by running the ``rack configure`` command. For more information, see :ref:`installation_and_configuration`.

You can get help for any command and its options by appending ``--help`` to the series of commands.

::

    $ rack servers instance create --help

Delete servers with an error status
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

You can list servers that have an error status and then delete them all in one line.

::

    $ rack servers instance list --status error --fields id --no-header | rack servers instance delete --stdin id

If you get a ``404 Not Found`` error message, it means that no servers were in error status.

Reboot multiple servers
~~~~~~~~~~~~~~~~~~~~~~~

With this one-line example, you list all the servers in an active state, omit the header row from the output, look for servers with ``-db`` in the name, and then restart those servers. Use the appropriate search command for your operating system. 

::

    $ rack servers instance list --status active --fields name --no-header |
    grep -i '-db' | rack servers instance reboot --soft --stdin name

Search for existing servers
~~~~~~~~~~~~~~~~~~~~~~~~~~~

If you have a large number of servers, you can use the following command to search through the list of available servers running in your account. Use the appropriate search command for your operating system. 

::

    $ rack servers instance list | grep "minecraft"

::

    543ce918-9d5c-476b-80a8-eefd396214ef	minecraft	ACTIVE	23.253.213.35	10.209.161.191	e19a734c-c7e6-443a-830c-242209c4d65d	performance1-4

If you expect a long list of servers in the output, you can list them with only the server ID returned.

::

    $ rack servers instance list --fields id

::

    ID
    aa049bf9-132c-4364-9808-bea21a009061
    543ce918-9d5c-476b-80a8-eefd396213ff

You can just get a list of IP addresses for all your cloud servers.

::

    $ rack servers instance list --fields publicipv4

::

    PublicIPv4
    162.209.0.32
    23.253.213.33

You can also search through metadata on each server. The following example shows the Orchestration information available for a particular server.

::

    $ rack servers instance get-metadata --name minecraft

::

    rax-heatdf149087-bf14-468c-9cfe-a76d83e43066

Get required information before creating a server
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

You usually want to list or create a few things before you launch  a server instance. For example, you might want to choose a flavor and image, and add a key pair that you can use to log in to the server after it is launched. This series of commands shows you what to expect in return.

If you want to upload your public key, you can either point to the file or copy and paste it into the command itself.

::

    $ cat ~/.ssh/id_rsa.pub

::

    ssh-rsa AAAB3.........t0mr name@example.com

    $ rack servers keypair upload --file ~/.ssh/id_rsa.pub --name macpub

or

::

    $ rack servers keypair upload --public-key "ssh-rsa AAAB3.........t0mr
    name@example.com" --name macpub

View any key pairs that you already have by listing them.

::

    $ rack servers keypair list

::

    Name                    Fingerprint
    4cb08c2f-c9db-4b00-86db-5d4b2c9a3aff    01:1b:4a:8f:9b:a3:c3:76:3d:90:06:bd:d2:5e:c2:16
    macpub                    5b:6e:55:2e:07:db:6c:e2:f6:4e:96:eb:29:30:64:2d

Now get the current list of images. Images created with Rackspace Cloud Images are listed first, followed by any snapshot images that you have stored in your account.

::

    $ rack servers image list | grep -i ubuntu

::

    973775ab-0653-4ef8-a571-7a2777787735	Ubuntu 12.04 LTS (Precise Pangolin) (PVHVM)		ACTIVE	20	512
    656e65f7-6441-46e8-978d-0d39beaaf559	Ubuntu 12.04 LTS (Precise Pangolin) (PV)		ACTIVE	20	512
    4315b2dc-23fc-4d81-9e73-aa620357e1d8	Ubuntu 15.04 (Vivid Vervet) (PVHVM)			ACTIVE	20	512
    09de0a66-3156-48b4-90a5-1cf25a905207	Ubuntu 14.04 LTS (Trusty Tahr) (PVHVM)			ACTIVE	20	512
    5ed162cc-b4eb-4371-b24a-a0ae73376c73	Ubuntu 14.04 LTS (Trusty Tahr) (PV)			ACTIVE	20	512

Next, choose the size and power of the server by looking at the available flavors.

::

    $ rack servers flavor list | grep -i compute

::

    compute1-15		15 GB Compute v1	15360	0	0	8	1250
    compute1-30		30 GB Compute v1	30720	0	0	16	2500
    compute1-4		3.75 GB Compute v1	3840	0	0	2	312.5
    compute1-60		60 GB Compute v1	61440	0	0	32	5000
    compute1-8		7.5 GB Compute v1	7680	0	0	4	625

Start a server with a key pair and metadata
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

In this example, you choose the image and flavor to launch a Rackspace Cloud Servers instance, such as a 4 GB General Purpose server on Ubuntu 14.04 LTS (Trusty Tahr) (PVHVM), and then put those values into the command along with the key pair and any metadata key-value pairs that you want to include. 

::

    $ rack servers instance create --name devserver \
    --image-name "Ubuntu 14.04 LTS (Trusty Tahr) (PVHVM)" \
    --flavor-id general1-4 --metadata purpose=dev \
    --keypair macpub

::

    ID        ab95d1d6-27d1-42bb-8cdc-800efcb5fc1f
    AdminPass    k6yfaDkgQfEr

Now you can view the server to ensure that the status is active. 

::

   $ rack servers instance list | grep devserver

::

    ID					Name		Status	PublicIPv4	PrivateIPv4Image					Flavor
    ab95d1d6-27d1-42bb-8cdc-800efcb5fc1f	devserver	ACTIVE	23.253.50.104	10.209.137.65	09de0a66-3156-48b4-90a5-1cf25a905207	general1-4

Connect to the server with SSH by using your public key.

::

    $ ssh root@23.253.50.104

Start a server from a volume
~~~~~~~~~~~~~~~~~~~~~~~~~~~~

The boot from volume feature gives you the ability to start a server with an attached volume. You can either start with a volume with a bootable image, to enable simpler migration when a server fails, or a storage volume that remains intact even after a server is shut down or deleted.

To create a bootable volume from an image and launch an instance from this volume, use the ``--block-device`` parameter. 

Use the following command to list your bootable volumes::

    $ rack block-storage volume list
    ID					Name		Bootable	Size	Status		VolumeType	SnapshotID
    18d361d1-2875-458b-9917-65010e37982a	BFV-test-SSD	true		100	in-use		SSD		
    88f2a1b0-b5f7-4634-ac4c-5e7ef0d9b2ac	BFB-test-SSD	true		100	available	SSD		
    6efa7008-ada7-4438-9033-efba4aa5cb06	Volume-1	false		100	available	SATA		

Following are the parameters for ``--block-device``:

- ``source-type=SOURCE``
    The type of object used to create the block device. Valid values are ``volume``, ``snapshot``, and ``image``.
    
- ``source-id=ID``
    The ID of the source resource (volume, snapshot, or image) from which to create the instance.

- ``destination-type=DEST``
    The type of the target virtual device. Valid values are ``volume`` and ``local``.

- ``volume-size=SIZE``
    The size of the volume that is created in GB.

- ``delete-on-termination={true|false}``
    What to do with the volume when the instance is deleted. Use ``false`` to delete the volume and use ``true`` to delete the
    volume when the instance is deleted. 

The following example command boots a server instance from a volume::

    $ rack servers instance create --name rackTestBFV  --block-device \
    "source-type=image,source-id=18d361d1-2875-458b-9917-65010e37982a,\
    volume-size=100,destination-type=volume,delete-on-termination=false" \
    --flavor-id compute1-15 --keypair macpub
