.. _servers:

=======
Servers
=======

This section lists the commands for interacting with Cloud Servers.

Commands
--------

All ``servers`` commands are based on this syntax::

   rack servers <subservice> <action> [command flags]

*Command flags* allow you to customize certain attributes of the command,
such as ``--name`` to name an instance. Type ``rack servers <subservice> <action> --help``
to bring up a list *command flags* specific to the command.


``instance``
~~~~~~~~~~~~

Instance commands use this syntax::

    rack servers instance <action> [command flags]

``list``
^^^^^^^^
Retrieves a list of virtual and bare metal servers::

    rack servers instance list [optional flags]

**Response**

.. code::

    $ rack servers instance list
    ID					                          Name		        Status	PublicIPv4	    PrivateIPv4	   Image					                       Flavor
    587c782d-2baa-493a-9545-f63bc6fa15ec	Cloud-Server-08	ACTIVE	104.130.252.96	10.209.10.28	 4315b2dc-23fc-4d81-9e73-aa620357e1d8	 general1-1
    00e54d18-06fd-4d6a-a7e9-900c7c80ebe7	serverBB	      ACTIVE	166.78.60.76	  10.182.7.74	   2f85a777-9ffd-4b49-a60e-1155ceb93a5e	 4
    c089d369-123a-4f6a-b06c-fa52f3218993	serverAA	      ACTIVE	162.209.72.242	10.182.1.134	 2f85a777-9ffd-4b49-a60e-1155ceb93a5e	 4
    c096178c-834a-41fa-a39c-f896d1abbe1b	serverB		      ACTIVE	166.78.138.250	10.182.14.215	 2f85a777-9ffd-4b49-a60e-1155ceb93a5e	 4
    e2f6b206-278d-40e4-915e-cce62a171ac0	ServerA		      ACTIVE	104.130.132.199	10.208.232.222 4315b2dc-23fc-4d81-9e73-aa620357e1d8	 general1-1
    8a254ea3-77b5-4f74-a893-8d2d51ae2cca	RACK		        ACTIVE	23.253.234.169	10.208.172.163						                           general1-1
    204408be-9c6c-41d8-8b59-265627089738	myUbuntuServer	ACTIVE	50.56.175.66	  10.180.2.239	 b66c7bf3-8677-4fe3-a117-aa6762acacec	 6
    33521191-cffe-435d-92f1-8bbd6d9f225a	Cloud-Server-06	ACTIVE	162.209.73.92 	10.182.1.155	 042395fc-728c-4763-86f9-9b0cacb00701	 2

``create``
^^^^^^^^^^
Creates a server instance::

    rack servers instance create --name <serverName> [optional flags]
    (echo serverName1 && echo serverName2) | rack servers instance create --stdin name [optional flags]

In order for this command to work, you must select an image using either the
`--image-id` or `--image-name` flags with the ID or name of the image you wish to use.
Alternatively, you may provide the `--block-device` flag to boot an instance from a volume.
In either case, you need to select a flavor by using the `--flavor-id` or `--flavor-name` with
the ID or name of the flavor you wish to use.

**Response**

.. code::

    $ rack servers instance create --name Rack4 --image-name "Red Hat Enterprise Linux 7 (PVHVM)" --flavor-id 2
    ID		e6a7263b-85ab-4640-b886-70eaaaf37e8c
    AdminPass	8n75vdF8EL2U

.. note::

    The admin password for your server is only given once. Please copy
    the password if you wish to make changes to server in the future.

``get``
^^^^^^^
Retrieves details of a specified server::

    rack servers instance get --id <serverID> [optional flags]
    rack servers instance get --name <serverName> [optional flags]
    (echo serverID1 && echo serverID2) | rack servers instance get --stdin id [optional flags]

**Response**

.. code::

    $ rack servers instance get --name Rack4
    ID					e6a7263b-85ab-4640-b886-70eaaaf37e8c
    Name					Rack4
    Status					ACTIVE
    Created					2015-08-06T18:12:59Z
    Updated					2015-08-06T18:16:21Z
    Image					bcc314ad-d971-4753-aea4-8b54d6219dfd
    Flavor					2
    PublicIPv4				166.78.63.84
    PublicIPv6				2001:4800:7812:514:be76:4eff:fe04:4f65
    PrivateIPv4				10.182.8.212
    KeyName
    TenantID				661145
    Progress				100
    SecurityGroups				[]
    Metadata:rax_service_level_automation	Complete

``update``
^^^^^^^^^^
Updates one or more editable attributes of a specified server
instance::

    rack servers instance update --id <serverID> [optional flags]
    rack servers instance update --name <serverName> [optional flags]

The response will reflect the updated attribute, depending on which attribute
you chose to update.

**Response**

.. code::

    $ rack servers instance update --rename ReRack
    ID		8a254ea3-77b5-4f74-a893-8d2d51ae2cca
    Name		ReRack
    PublicIPv4	23.253.234.169
    PublicIPv6	2001:4800:7817:103:be76:4eff:fe04:853f


``delete``
^^^^^^^^^^
Deletes a server instance::

    rack servers instance delete --id <serverID> [optional flags]
    rack servers instance delete --name <serverName> [optional flags]
    (echo serverID1 && echo serverID2) | rack servers instance delete --stdin id [optional flags]

**Response**

.. code::

    $ rack servers instance delete --name ReRack
    Deleting instance [8a254ea3-77b5-4f74-a893-8d2d51ae2cca]

``reboot``
^^^^^^^^^^
Performs a soft or hard reboot of a specified server. A soft reboot
will slowly shutdown and restart your server's operating system. A hard reboot
will perform an immediate shutdown and restart::

    rack servers instance reboot --id <serverID> [--soft | --hard] [optional flags]
    rack servers instance reboot --name <serverName> [--soft | --hard] [optional flags]
    (echo serverID1 && echo serverID2) | rack servers instance reboot --stdin id [--soft | --hard] [optional flags]

*Response*

.. code::

    $ rack servers instance reboot --name RACK --hard
    Successfully rebooted instance [0807eefe-b36a-415c-bd59-8b4cef63c563]

``rebuild``
^^^^^^^^^^^
Removes all data on the server and replaces it with the specified image::

    rack servers instance rebuild --id <serverID> --image-id <imageID> --admin-pass <adminPass> [optional flags]
    rack servers instance rebuild --name <serverName> --image-id <imageID> --admin-pass <adminPass> [optional flags]

**Response**

.. code::

    $ rack servers instance rebuild --name Rack4 --image-id ab5c119f-50ab-4213-b969-19b1853d41b0 --admin-pass 8n75vdF8EL2U
    Successfully rebuilt instance [0807eefe-b36a-415c-bd59-8b4cef63c563]

``resize``
^^^^^^^^^^
Converts an existing server to a different flavor, which scales the
server up or down. The original server is saved for a period of time to allow roll
back if a problem occurs::

    rack servers instance resize --id <serverID> --flavor-id <flavorID> [optional flags]
    rack servers instance resize --name <serverName> --flavor-id <flavorID> [optional flags]
    (echo serverID1 && echo serverID2) | rack servers instance resize --stdin id --flavor-id <flavorID> [optional flags]

**Response**

.. code:

    $ rack servers instance resize --name Rack4 --flavor-id 4
    Successfully resized instance [e6a7263b-85ab-4640-b886-70eaaaf37e8c] to flavor [4]

.. note::
    This command is not available for OnMetal servers.

``set-metadata``
^^^^^^^^^^^^^^^^
Sets metadata for the specified server or image::

    rack servers instance set-metadata --id <serverID> --metadata <key1=val1,key2=val2,...> [optional flags]
    rack servers instance set-metadata --name <serverName> --metadata <key1=val1,key2=val2,...> [optional flags]

**Response**

.. code::

    $ rack servers instance set-metadata --id e6a7263b-85ab-4640-b886-70eaaaf37e8c --metadata heat=true
    Metadata:heat	true

``get-metadata``
^^^^^^^^^^^^^^^^
Retrieves the metadata for a given server::

    rack servers instance get-metadata --id <serverID> [optional flags]
    rack servers instance get-metadata --name <serverName> [optional flags]

**Response**

.. code::

    $ rack servers instance get-metadata --name Rack4
    Metadata:heat	true

``update-metadata``
^^^^^^^^^^^^^^^^^^^
Updates metadata items for a specified server or image, or adds the specified
metadata if there is no current metadata associated with the server or image::

    rack servers instance update-metadata --id <serverID> --metadata <key1=val1,key2=val2,...> [optional flags]
    rack servers instance update-metadata --name <serverName> --metadata <key1=val1,key2=val2,...> [optional flags]

**Response**

.. code::

    $ rack servers instance update-metadata --name Rack4 --metadata heat=false
    Metadata:heat	false


``delete-metadata``
^^^^^^^^^^^^^^^^^^^
Deletes one or more metadata keys from a server::

    rack servers instance delete-metadata --id <serverID> --metadata-keys <key1,key2,...> [optional flags]
    rack servers instance delete-metadata --name <serverName> --metadata-keys <key1,key2,...> [optional flags]

**Response**

.. code::

    $ rack servers instance delete-metadata --name Rack4 --metadata-keys heat
    Successfully deleted metadata

``image``
~~~~~~~~~

Image commands use this syntax::

    rack server image <action> [optional flags]

``list``
^^^^^^^^
Lists all images visible by your account::

    rack servers image list [optional flags]

**Response**

.. code::

    $ rack servers image list
    ID					                          Name							                     	    Status	MinDisk	MinRAM
    faad95b7-396d-483e-b4ae-77afec7e7097	Vyatta Network OS 6.7R9					            ACTIVE	20	    1024
    5a2a568b-0a56-4821-82b5-279bbca7cd9d	Windows Server 2012 R2						          ACTIVE	40	    1024
    c934d497-7b45-4764-ac63-5b67e1458a20	Debian 7 (Wheezy) (PVHVM)					          ACTIVE	20	    512
    973775ab-0653-4ef8-a571-7a2777787735	Ubuntu 12.04 LTS (Precise Pangolin) (PVHVM)	ACTIVE	20	    512
    656e65f7-6441-46e8-978d-0d39beaaf559	Ubuntu 12.04 LTS (Precise Pangolin) (PV)		ACTIVE	20	    512
    2cc5db1b-2fc8-42ae-8afb-d30c68037f02	Fedora 22 (PVHVM)						                ACTIVE	20	    512
    c25f1ae0-30b3-4012-8ca6-5ecfcf05c965	CentOS 7 (PVHVM)						                ACTIVE	20	    512
    3cdcd2cc-238c-4f42-a9f4-0a80de217f7a	OpenSUSE 13.2 (PVHVM)					            	ACTIVE	20	    512
    36076d08-3e8b-4436-9253-7a8868e4f4d7	Scientific Linux 6 (PVHVM)					        ACTIVE	20	    512
    ab5c119f-50ab-4213-b969-19b1853d41b0	Scientific Linux 7 (PVHVM)					        ACTIVE	20	    512
    7a1cf8de-7721-4d56-900b-1e65def2ada5	FreeBSD 10 (PVHVM)						              ACTIVE	20	    512
    168c1be2-a3b0-423f-a619-f63cce550063	Gentoo 15.3 (PVHVM)						              ACTIVE	20	    512
    4315b2dc-23fc-4d81-9e73-aa620357e1d8	Ubuntu 15.04 (Vivid Vervet) (PVHVM)				  ACTIVE	20	    512
    ade87903-9d82-4584-9cc1-204870011de0	Arch 2015.7 (PVHVM)						              ACTIVE	20	    512

``get``
^^^^^^^
Retrieves details of the specified image::

    rack servers image get --id <imageID> [optional flags]
    rack servers image get --name <imageName>] [optional flags]
    (echo imageID1 && echo imageID2) | rack servers image get --stdin id [optional flags]

**Response**

.. code::

    $ rack servers image get --id bcc314ad-d971-4753-aea4-8b54d6219dfd
    ID	bcc314ad-d971-4753-aea4-8b54d6219dfd
    Name	Red Hat Enterprise Linux 7 (PVHVM)
    Status	ACTIVE
    Progress100
    MinDisk	20
    MinRAM	512
    Created	2015-07-27T17:57:55Z
    Updated	2015-07-28T20:34:24Z

.. note::

   To guarantee usage of the same image every time, use the `--id` flag. Images often
   are updated with security patches, and the updated images will have a different ID but
   the same name.


``flavor``
~~~~~~~~~~

Flavor commands use this syntax::

    rack servers flavor <action> [optional flags]

``list``
^^^^^^^^
Lists information for all available flavors::

    rack servers flavor list [optional flags]

**Response**

.. code::

    $ rack servers flavor list
    ID			Name			              RAM	  Disk	Swap	VCPUs	RxTxFactor
    2			  512MB Standard Instance	512	  20	  512	  1	    80
    3			  1GB Standard Instance	  1024	40	  1024	1	    120
    4			  2GB Standard Instance	  2048	80	  2048	2	    240
    5			  4GB Standard Instance	  4096	160	  2048	2	    400
    6			  8GB Standard Instance	  8192	320	  2048	4	    600
    7			  15GB Standard Instance	15360	620	  2048	6	    800
    8			  30GB Standard Instance	30720	1200	2048	8	    1200

``get``
^^^^^^^
Retrieves details of the specified flavor::

    rack servers flavor get --id <flavorID> [optional flags]
    rack servers flavor get --name <flavorName>] [optional flags]
    (echo flavorID1 && echo flavorID2) | rack servers flavor get --stdin id [optional flags]

**Response**

.. code::

    $ rack servers flavor get --id 4
    ID			                 4
    Name			               2GB Standard Instance
    Disk			               80
    RAM			                 2048
    RxTxFactor		           240
    Swap			               2048
    VCPUs			               2
    ExtraSpecs:PolicyClass	 standard_flavor
    ExtraSpecs:NumDataDisks	 0
    ExtraSpecs:Class	       standard1
    ExtraSpecs:DiskIOIndex	 0

``keypair``
~~~~~~~~~~~

Keypair commands use this syntax::

    rack servers keypair <action> [optional flags]

``list``
^^^^^^^^
Retrieves a list of all key pairs associated with your account::

    rack servers keypair list [flags]

**Response**

.. code::

    Name					                        Fingerprint
    770fb26f-2c43-4196-95d1-ad9ec1008365	1d:86:3b:a4:19:d9:01:f8:16:83:d3:43:6a:10:98:87
    public key				                    7e:e6:8e:af:64:5b:d7:aa:4c:9c:ea:c8:05:0a:29:2e

``generate``
^^^^^^^^^^^^
Generates a newly created key pair with the specified name::

    rack servers keypair generate --name <keypairName> [optional flags]
    (echo keypairName1 && echo keypairName2) | rack servers keypair generate --stdin name [optional flags]

**Response**

.. code::

    $ rack servers keypair generate --name "rack key"
    Name		    rack key
    Fingerprint	73:5d:f5:1d:2d:00:29:59:4c:82:66:f4:10:58:c3:7e
    PublicKey	  ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCdhmClpS/NF8RGf9Lhj9ffmOm7iUeZd0Mr7CmS+tHwfLLrFfb5VWIQm8E5vnkGbA4iCE1XDC6BjTqcdhsfJtRoyw0HRGcdbHNM2r5muXSdl0r1aRz5jzPUH1e0Ot88UW3YRO8rWAqzUVbRIm2T/K24g8HSs8qDQCMrS4m/tJH4sKKRVhr/CYGs7mYflCh0Y7bHOuJCkMpCWSV4e+2xHciIqgpDS7aduQAo8pFdza6lV9n1QiJ2sSJnoH9IKVzw9RfJNsVS9hsqMB/GFKIrDtmABYcbuDQ0OXrZQusF/hbzXXJc89uRdG2/aP7NUJfSvcLCJXxxoixMddcZOFOjEz8H Generated-by-Nova
    PrivateKey	-----BEGIN RSA PRIVATE KEY-----
                MIIEqAIBAAKCAQEAnYZgpaUvzRfERn/S4Y/X35jpu4lHmXdDK+wpkvrR8Hyy6xX2
                +VViEJvBOb55BmwOIghNVwwugY06nHYbHybUaMsNB0RnHWxzTNq+Zrl0nZdK9Wkc
                +Y8z1B9XtDrfPFFt2ETvK1gKs1FW0SJtk/ytuIPB0rPKg0AjK0uJv7SR+LCikVYa
                /wmBrO5mH5QodGO2xzriQpDKQlkleHvtsR3IiKoKQ0u2nbkAKPKRXc2upVfZ9UIi
                drEiZ6B/SClc8PUXyTbFUvYbKjAfxhSiKw7ZgAWHG7g0NDl62ULrBf4W811yXPPb
                kXRtv2j+zVCX0r3CwiV8caIsTHXXGThToxM/BwIDAQABAoIBAF9U6byVob4vOfuo
                dAlKMk0Bj8KNLCW2RZKZU/e9j7xI20bmfhcbE0QV3vSBT2ERl+QEvjmGB8cjn2r9
                FyDdMQNSj0nsnbLy5TyHzV8BsY+F9jlbKcvmXGltpnhOgLdSWHXgEhZbK+6ltTLP
                8Rz0eHYIVw6a/F4+PIqBJpY8Q3MLD9To6+Nu+ytcnoRpW13ItyTkg61Q60tABMSP
                wHHWkMYMhxnCqvJA+0n9Hkj4l+xZffbPKQkAd/oFbO+/3qwUys28zY5pdHhO0RzD
                vQabzw+UxsMIOe3IwaP78fP2Sw5AV6ruuQ7mGCfZxLq8Of0AyjVL7Adoj6T9WzvO
                FBLp2nECggCBAMP9Mlb1zp+dtfVrohw4gVgAqs9W4K+Eh6FstRqZnM8Iq4sSpoyf
                rzeRVeRw4hEiQdyDGyUBgkFCB2zIeB3FpuVdGN3j13A6xD5J2wpXamn1ysdE3ShA
                tJWozZFK82t5wTnuWGLwoZyNHKaiXiUXQUXlpTYJcXhrHLbqQsPCyrMLAoIAgQDN
                wiFBmbPteAfZ6w4gYKyzgB7BGp8ZaSb1/Z/80Q1r5XN8AzTqPciiZqHgRvfsve3k
                c7UA/mjs4MiSTdURFYS00EScKpDoSyu9Q9vCIKcLo/ijxrMz+3QrN14QHFlHkXG8
                t1JXYHop6HVxdBIiovDreovOpA8KuPmY2ZvbivEhdQKCAIBtvWBqgMhpZ41zFgl3
                c8w40mbSIfs9TCAFqzDc7BZ5dolXHeJT6zXZD2/gsbRjk6L9LgwM9INStv4hUz6u
                rPV+iFpcJC+Fj+JCXmhfqgLTweUBHvYWi+SRyCsSp70U+N/Q6DxlageT+Q+J3nCI
                pDTQRn3ze+YlXxR89z7bDj8hcQKCAIEAxUQJfktOJa2eWV7x/DizWqTK13gecM5P
                fCfc6xXCOF+TiKHKaYkyggDD0bI4n9C38v672mgPUItxwFK+o9JtmKzUGqT0qMDo
                /lvApS2I+bALAXnO9Vdu2MAMfoWvUt4unS9k0kC83tDvSAZwHKT7NcgXodXIVg9h
                vRlkQ+fBpsECggCADIUPZDRtqFiBnKYI1sywCAT50plRs7o0yRcFtJyp4rQczLbO
                6fYay1fgBrYW8CxHnJfeP/zCFGGxDxYjnbqI3GKGQHqFCegkxirAx7gEM6sllG4g
                EywgWCyPegcAZe1TjH3VfAr+nroMpURJKB6YMjdyh/o7xkm/NaC2cbNR6jc=
                -----END RSA PRIVATE KEY-----

``upload``
^^^^^^^^^^
Uploads an existing key pair with the specified name::

    rack servers keypair upload --name <keypairName> --public-key <publicKeyData> [optional flags]
    rack servers keypair upload --name <keypairName> --file <publicKeyfile> [optional flags]

**Response**

.. code::

    $ rack servers keypair upload --name racksa --public-key "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDY31xd7OV8vrEYseCRTuEOfGVJRpVRIUdyQT21mp1nfbXV4vSWj2cOsi6kv+HsbxHaAok6LCRA8YUzGqLcQG+5cAUHQ7GPHsaWrTlP/OmcN7BOxFSazGEldQWjm04EW6ahRXrMthrg7L84O4b7RgNA9REmHkhqb5dqXcuIb99fRh/9f6cUIzqyduB9pcmSafY0uzRrUmrkjmSwj1dIifKlsfjHku8RChjBQGTMc+3c6Cjr+TlnvMoBOoemK1kxb0HJDFJZtXdN8RvVwzvLD7EZfBZZ4exew9u+hWpV0G2H8jDQQeHDErTsIUhWVZQxFgR8uknGWXt/du7Y4d0NJ7GP nath8916@MPM1XEDV30"
    Name		racksa
    Fingerprint	5d:2c:fe:90:fc:42:89:70:d1:7d:2e:ad:a1:31:a8:a2
    PublicKey	ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDY31xd7OV8vrEYseCRTuEOfGVJRpVRIUdyQT21mp1nfbXV4vSWj2cOsi6kv+HsbxHaAok6LCRA8YUzGqLcQG+5cAUHQ7GPHsaWrTlP/OmcN7BOxFSazGEldQWjm04EW6ahRXrMthrg7L84O4b7RgNA9REmHkhqb5dqXcuIb99fRh/9f6cUIzqyduB9pcmSafY0uzRrUmrkjmSwj1dIifKlsfjHku8RChjBQGTMc+3c6Cjr+TlnvMoBOoemK1kxb0HJDFJZtXdN8RvVwzvLD7EZfBZZ4exew9u+hWpV0G2H8jDQQeHDErTsIUhWVZQxFgR8uknGWXt/du7Y4d0NJ7GP nath8916@MPM1XEDV30
    PrivateKey

``get``
^^^^^^^
Retrieves details on a specified key pair::

    rack servers keypair get --name <keypairName> [optional flags]
    (echo keypairName1 && echo keypairName2) | rack servers keypair get --stdin name [optional flags]

**Response*

.. code::

    Name		    rack key
    Fingerprint	73:5d:f5:1d:2d:00:29:59:4c:82:66:f4:10:58:c3:7e
    PublicKey	  ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCdhmClpS/NF8RGf9Lhj9ffmOm7iUeZd0Mr7CmS+tHwfLLrFfb5VWIQm8E5vnkGbA4iCE1XDC6BjTqcdhsfJtRoyw0HRGcdbHNM2r5muXSdl0r1aRz5jzPUH1e0Ot88UW3YRO8rWAqzUVbRIm2T/K24g8HSs8qDQCMrS4m/tJH4sKKRVhr/CYGs7mYflCh0Y7bHOuJCkMpCWSV4e+2xHciIqgpDS7aduQAo8pFdza6lV9n1QiJ2sSJnoH9IKVzw9RfJNsVS9hsqMB/GFKIrDtmABYcbuDQ0OXrZQusF/hbzXXJc89uRdG2/aP7NUJfSvcLCJXxxoixMddcZOFOjEz8H Generated-by-Nova
    UserID		  172157

``delete``
^^^^^^^^^^
Deletes the specified key paid::

    rack servers keypair delete --name <keypairName> [optional flags]
    (echo keypairName1 && echo keypairName2) | rack servers keypair delete --stdin name [optional flags]

**Response**

.. code::

    $ rack servers keypair delete --name "rack key"
    Successfully deleted keypair [rack key]


``volume-attachment``
~~~~~~~~~~~~~~~~~~~~~

Volume attachment commands use this syntax::

    rack server volume-attachment <action> [optional flag]

These commands are often used with :ref:`Cloud Block Storage <blockexamples>`.

``list``
^^^^^^^^
Lists the volume attachments for the specified server::

    rack servers volume-attachment list --server-id <serverID> [optional flags]
    rack servers volume-attachment list --server-name <serverName> [optional flags]
    rack servers volume-attachment list --stdin server-id [optional flags]

**Response**

.. code::

    $ rack servers volume-attachment list --server-name Rack4
    ID					                          Device		VolumeID				                      ServerID
    d823ddd9-73dc-426e-8d4c-38841941dd57	/dev/xvdb	d823ddd9-73dc-426e-8d4c-38841941dd57	e6a7263b-85ab-4640-b886-70eaaaf37e8c
    8349b7c7-acf0-4c5f-9bae-38fc87d0142d	/dev/xvdd	8349b7c7-acf0-4c5f-9bae-38fc87d0142d	e6a7263b-85ab-4640-b886-70eaaaf37e8c

``create``
^^^^^^^^^^
Attaches one or more volumes to the specified sever::

    rack servers volume-attachment create --server-id <serverID> --volume-id <volumeID> [optional flags]
    rack servers volume-attachment create --server-name <serverName> --volume-id <volumeID> [optional flags]
    rack servers volume-attachment create --server-id <serverID> --volume-name <volumeName> [optional flags]
    rack servers volume-attachment create --server-name <serverName> --volume-name <volumeName> [optional flags]
    (echo volumeID1 && echo volumeID2) | rack servers volume-attachment create --server-id <serverID> --stdin volume-id [optional flags]
    (echo volumeID1 && echo volumeID2) | rack servers volume-attachment create --server-name <serverName> --stdin volume-id [optional flags]

**Response**

.. code::
    $ rack servers volume-attachment create --server-name Rack4 --volume-id 8349b7c7-acf0-4c5f-9bae-38fc87d0142d
    ID	8349b7c7-acf0-4c5f-9bae-38fc87d0142d
    Device	/dev/xvdd
    VolumeID8349b7c7-acf0-4c5f-9bae-38fc87d0142d
    ServerIDe6a7263b-85ab-4640-b886-70eaaaf37e8c

``get``
^^^^^^^
Retrieves details of a specified volume attachment ID for a specified server::

    rack servers volume-attachment get --server-id <serverID> --id <attachmentID> [optional flags]
    rack servers volume-attachment get --server-name <serverName> --id <attachmentID> [optional flags]

**Response**

.. code::

    $ rack servers volume-attachment get --server-name Rack4 --id d823ddd9-73dc-426e-8d4c-38841941dd57
    ID	d823ddd9-73dc-426e-8d4c-38841941dd57
    Device	/dev/xvdb
    VolumeIDd823ddd9-73dc-426e-8d4c-38841941dd57
    ServerIDe6a7263b-85ab-4640-b886-70eaaaf37e8c

``delete``
^^^^^^^^^^
Removes a specified volume attachment from a specified server instance::

    rack servers volume-attachment delete --server-id <serverID> --id <attachmentID> [optional flags]
    rack servers volume-attachment delete --server-name <serverName> --id <attachmentID> [optional flags]

**Response**

.. code::

    $ rack servers volume-attachment delete --server-name Rack4 --id d823ddd9-73dc-426e-8d4c-38841941dd57
    Successfully deleted volume attachment [d823ddd9-73dc-426e-8d4c-38841941dd57]
