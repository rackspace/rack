.. _serversexamples:

======================
Cloud Servers examples
======================

Before you get started on any examples, be sure you have entered your
username and API key and stored them locally::

    rack configure

You can get help for any command and its options by appending --help to the
series of commands::

    $ rack servers instance create --help

Delete servers with ERROR status
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

You can list servers in ERROR status and then delete them all in one line::

    $ rack servers instance list --status error --fields id | tail -n+2 | rack servers instance delete --stdin id

If you get a 404 Not Found, it means no servers were in error status.

Reboot multiple servers
~~~~~~~~~~~~~~~~~~~~~~~

With this one-line example, you list all the server IDs in an ACTIVE state, cut
the header row that's output, then look for your servers with "-db" in the name
and restart them.::

    $ rack servers instance list --status active --fields name --no-header |
    grep -i '-db' | rack servers instance reboot --soft --stdin name

Search for existing servers
~~~~~~~~~~~~~~~~~~~~~~~~~~~

If you have a lot of servers, the `rack` command lets you search through
the list. 

Use grep to search through the list of available Cloud Servers running in your
account::

    $ rack servers instance list | grep "minecraft"

::

    543ce918-9d5c-476b-80a8-eefd396214ef	minecraft	ACTIVE	23.253.213.35	10.209.161.191	e19a734c-c7e6-443a-830c-242209c4d65d	performance1-4

If you have a long list of servers, here's an example of listing with only the
server ID returned::

    $ rack servers instance list --fields id

::

    ID
    aa049bf9-132c-4364-9808-bea21a009061
    543ce918-9d5c-476b-80a8-eefd396213ff

Or just get a list of IP addresses for all your Rackspace Cloud Servers::

    $ rack servers instance list --fields publicipv4

::

    PublicIPv4
    162.209.0.32
    23.253.213.33

Or search through metadata on each server. This example shows the Orchestration
information available on this particular server::

    $ rack servers instance get-metadata --name minecraft

::

    rax-heatdf149087-bf14-468c-9cfe-a76d83e43066

Get required information before creating a server
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

You want to list or create a few things prior to launching a server instance,
such as choosing a flavor and image and also adding a keypair you can use to
log into the server once it is launched. This series of commands show you what
to expect in return.

If you want to upload your public key, you can either point to the file or
copy and paste it into the command itself::

    $ cat ~/.ssh/id_rsa.pub

::

    ssh-rsa AAAB3.........t0mr name@example.com

    $ rack servers keypair upload --file ~/.ssh/id_rsa.pub --name macpub

or::

    $ rack servers keypair upload --public-key ssh-rsa AAAB3.........t0mr
    name@example.com --name macpub

Take a look at any keypairs you already have by listing them::

    $ rack servers keypair list

::

    Name                    Fingerprint
    4cb08c2f-c9db-4b00-86db-5d4b2c9a3aff    01:1b:4a:8f:9b:a3:c3:76:3d:90:06:bd:d2:5e:c2:16
    macpub                    5b:6e:55:2e:07:db:6c:e2:f6:4e:96:eb:29:30:64:2d

Now get the current list of images. First are the Rackspace Cloud Images
followed by any snapshot images you have stored in your account::

    $ rack servers image list | grep -i ubuntu

::

    973775ab-0653-4ef8-a571-7a2777787735	Ubuntu 12.04 LTS (Precise Pangolin) (PVHVM)		ACTIVE	20	512
    656e65f7-6441-46e8-978d-0d39beaaf559	Ubuntu 12.04 LTS (Precise Pangolin) (PV)		ACTIVE	20	512
    4315b2dc-23fc-4d81-9e73-aa620357e1d8	Ubuntu 15.04 (Vivid Vervet) (PVHVM)			ACTIVE	20	512
    09de0a66-3156-48b4-90a5-1cf25a905207	Ubuntu 14.04 LTS (Trusty Tahr) (PVHVM)			ACTIVE	20	512
    5ed162cc-b4eb-4371-b24a-a0ae73376c73	Ubuntu 14.04 LTS (Trusty Tahr) (PV)			ACTIVE	20	512

Next, choose the size and power of the server by looking at the available
flavors::

    $ rack servers flavor list | grep -i compute

::

    compute1-15		15 GB Compute v1	15360	0	0	8	1250
    compute1-30		30 GB Compute v1	30720	0	0	16	2500
    compute1-4		3.75 GB Compute v1	3840	0	0	2	312.5
    compute1-60		60 GB Compute v1	61440	0	0	32	5000
    compute1-8		7.5 GB Compute v1	7680	0	0	4	625

Start a server with a keypair and metadata
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Choose the image and flavor to launch a Rackspace Cloud Server, such as
a 4 GB General Purpose on Ubuntu 14.04 LTS (Trusty Tahr) (PVHVM) and put those
values into the command along with the keypair and any metadata key-value pairs
you want to include. Here is an example::

    $ rack servers instance create --name devserver \
    --image-name "Ubuntu 14.04 LTS (Trusty Tahr) (PVHVM)" \
    --flavor-id general1-4 --metadata purpose=dev \
    --keypair macpub

::

    ID        ab95d1d6-27d1-42bb-8cdc-800efcb5fc1f
    AdminPass    k6yfaDkgQfEr

Now you can view the server to make sure the Status is ACTIVE::

   $ rack servers instance list | grep devserver

::

    ID					Name		Status	PublicIPv4	PrivateIPv4Image					Flavor
    ab95d1d6-27d1-42bb-8cdc-800efcb5fc1f	devserver	ACTIVE	23.253.50.104	10.209.137.65	09de0a66-3156-48b4-90a5-1cf25a905207	general1-4

To connect to the server with SSH using your public key, use this command::

    $ ssh root@23.253.50.104

Start a server from a volume
~~~~~~~~~~~~~~~~~~~~~~~~~~~~

The boot from volume features gives you the ability to start a server with an
attached volume. You can either start with a volume with a bootable image, to
enable simpler migration when a server fails, or a storage volume that remains
intact even after a server is shutdown or deleted.

To create a bootable volume from an image and launch an instance from
this volume, use the ``--block-device`` parameter.

   The settings are:

   -  ``--block-device``
      source=SOURCE,id=ID,dest=DEST,size=SIZE,shutdown=PRESERVE,bootindex=INDEX

         **source-type=SOURCE**
             The type of object used to create the block device. Valid values
             are ``volume``, ``snapshot``, ``image``, and ``blank``.

         **source-id=ID**
             The ID of the source. Use a volume ID if the ``source-type`` is
             a volume and an image ID if the ``source-type`` is image.

         **destination-type=DEST**
             The type of the target virtual device. Valid values are ``volume``
             and ``local``.

         **volume-size=SIZE**
             The size of the volume that is created in GB.

         **delete-on-termination={true\|false}**
             What to do with the volume when the instance is deleted. Use
             ``false`` to delete the volume and ``true`` to delete the
             volume when the instance is deleted.

Use this command to boot from a volume::

    $ rack servers instance create --name rackTestBFV  --block-device \
    "source-type=image,source-id=18d361d1-2875-458b-9917-65010e37982a,\
    volume-size=100,destination-type=volume,delete-on-termination=false" \
    --flavor-id compute1-15 --keypair macpub