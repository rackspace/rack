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

Search for existing servers
~~~~~~~~~~~~~~~~~~~~~~~~~~~

If you have a lot of servers, the `rack` command lets you search through
the list. 

Use grep to search through the list of available Cloud Servers running in your
account::

    $ rack servers instance list | grep "minecraft"

ID                    Name        Status    PublicIPv4    PrivateIPv4    Image                    Flavor
543ce918-9d5c-476b-80a8-eefd396214ff    minecraft    ACTIVE    23.253.213.33    10.209.161.191    e19a734c-c7e6-443a-830c-242209c4d65c    performance1-4

If you have a long list of servers, here's an example of listing with only the
server ID returned::

    $ rack servers instance list --fields id

    ID
    aa049bf9-132c-4364-9808-bea21a009061
    543ce918-9d5c-476b-80a8-eefd396213ff

Or just get a list of IP addresses for all your Rackspace Cloud Servers::

    $ rack servers instance list --fields publicipv4

    PublicIPv4
    162.209.0.32
    23.253.213.33

Or search through metadata on each server. This example shows the Orchestration
information available on this particular server::

    $ rack servers instance get-metadata --name minecraft
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

    ssh-rsa AAAB3.........t0mr name@example.com

    $ rack servers keypair upload --file ~/.ssh/id_rsa.pub --name macpub

or

    $ rack servers keypair upload --public-key ssh-rsa AAAB3.........t0mr name@example.com --name macpub

Take a look at any keypairs you already have by listing them::

    $ rack servers keypair list

    Name                    Fingerprint
    4cb08c2f-c9db-4b00-86db-5d4b2c9a3aff    01:1b:4a:8f:9b:a3:c3:76:3d:90:06:bd:d2:5e:c2:16
    macpub                    5b:6e:55:2e:07:db:6c:e2:f6:4e:96:eb:29:30:64:2d

Now get the current list of images. First are the Rackspace Cloud Images
followed by any snapshot images you have stored in your account::

    $ rack servers image list

    ID                    Name                                Status    MinDisk    MinRAM
    973775ab-0653-4ef8-a571-7a2777787735    Ubuntu 12.04 LTS (Precise Pangolin) (PVHVM)            ACTIVE    20    512
    656e65f7-6441-46e8-978d-0d39beaaf559    Ubuntu 12.04 LTS (Precise Pangolin) (PV)            ACTIVE    20    512
    2cc5db1b-2fc8-42ae-8afb-d30c68037f02    Fedora 22 (PVHVM)                        ACTIVE    20    512
    c25f1ae0-30b3-4012-8ca6-5ecfcf05c965    CentOS 7 (PVHVM)                        ACTIVE    20    512
    3cdcd2cc-238c-4f42-a9f4-0a80de217f7a    OpenSUSE 13.2 (PVHVM)                        ACTIVE    20    512
    36076d08-3e8b-4436-9253-7a8868e4f4d7    Scientific Linux 6 (PVHVM)                    ACTIVE    20    512
    ab5c119f-50ab-4213-b969-19b1853d41b0    Scientific Linux 7 (PVHVM)                    ACTIVE    20    512
    7a1cf8de-7721-4d56-900b-1e65def2ada5    FreeBSD 10 (PVHVM)                        ACTIVE    20    512
    5b7fdeda-c283-4797-a55c-e5d8411c51b4    Debian 7 (Wheezy) (PVHVM)                    ACTIVE    20    512
    168c1be2-a3b0-423f-a619-f63cce550063    Gentoo 15.3 (PVHVM)                        ACTIVE    20    512
    4315b2dc-23fc-4d81-9e73-aa620357e1d8    Ubuntu 15.04 (Vivid Vervet) (PVHVM)                ACTIVE    20    512
    ade87903-9d82-4584-9cc1-204870011de0    Arch 2015.7 (PVHVM)                        ACTIVE    20    512
    bcc314ad-d971-4753-aea4-8b54d6219dfd    Red Hat Enterprise Linux 7 (PVHVM)                ACTIVE    20    512
    8e3d8c5b-ac07-429f-8304-d2863e1a0636    Red Hat Enterprise Linux 6 (PVHVM)                ACTIVE    20    512
    783f71f4-d2d8-4d38-b2e1-8c916de79a38    Red Hat Enterprise Linux 6 (PV)                    ACTIVE    20    512
    05dd965d-84ce-451b-9ca1-83a134e523c3    Red Hat Enterprise Linux 5 (PV)                    ACTIVE    20    512
    6c21b351-e12a-4ddf-a0a0-a6849c2b0037    Fedora 21 (PVHVM)                        ACTIVE    20    512
    498c59a0-3c26-4357-92c0-dd938baca3db    Debian Unstable (Sid) (PVHVM)                    ACTIVE    20    512
    0535b46b-fae2-4813-945f-701949c53c2e    Debian Testing (Stretch) (PVHVM)                ACTIVE    20    512
    19149d8b-bd6a-4b0b-a688-657780f9cf6c    Debian 8 (Jessie) (PVHVM)                    ACTIVE    20    512
    09de0a66-3156-48b4-90a5-1cf25a905207    Ubuntu 14.04 LTS (Trusty Tahr) (PVHVM)                ACTIVE    20    512
    5ed162cc-b4eb-4371-b24a-a0ae73376c73    Ubuntu 14.04 LTS (Trusty Tahr) (PV)                ACTIVE    20    512
    aa68fd54-2f9a-42c3-9901-4035e2738830    CentOS 6 (PVHVM)                        ACTIVE    20    512
    21612eaf-a350-4047-b06f-6bb8a8a7bd99    CentOS 6 (PV)                            ACTIVE    20    512
    d75bc322-b02c-493d-b414-097b3bcce4dd    CentOS 5 (PV)                            ACTIVE    20    512
    0c934c81-6647-4166-921d-2250e8c6ff29    CoreOS (Stable)                            ACTIVE    20    512
    7cf9f618-3dd3-4e3e-bace-e44d857039e2    CoreOS (Beta)                            ACTIVE    20    512
    40c92b20-8b73-4732-9ea9-da47562cebe1    CoreOS (Alpha)                            ACTIVE    20    512
    2f572661-00f0-47a3-8955-c489e8269ed0    Vyatta Network OS 6.7R6                        ACTIVE    20    1024
    d77bb4fc-23ba-402b-af8f-1dfd0ba5df48    Windows Server 2012 R2 + SQL Server 2014 Web            ACTIVE    40    2048
    0dbfedcc-36ec-4425-bb35-6e0689b21475    Windows Server 2008 R2 SP1 + SQL Server 2008 R2 SP2 Web        ACTIVE    40    2048
    eed7d3f8-a8ec-4fd8-93e1-c3d87e826297    Windows Server 2012 R2 + SQL Server 2014 Standard        ACTIVE    40    2048
    3703dd75-0c02-46c4-9511-6fae6d4d9e25    Windows Server 2012 + SQL Server 2012 SP1 Standard        ACTIVE    40    2048
    a7e39b24-3aca-4a96-8b7c-aba3c720d54d    Windows Server 2012 R2                        ACTIVE    40    1024
    a512248f-c5a7-4159-ab20-8762d5c8c93e    Windows Server 2008 R2 SP1 + SQL Server 2008 R2 SP2 Standard    ACTIVE    40    2048
    9f65c6e3-3253-46f3-92d1-5d2681282590    Windows Server 2012 + SQL Server 2012 SP1 Web            ACTIVE    40    2048
    526c377d-d29a-41fd-9227-1cecd2bf418d    Windows Server 2008 R2 SP1                    ACTIVE    40    1024
    1229d02f-5189-4dfe-a255-7d285a2f0bc9    Windows Server 2012                        ACTIVE    40    1024

Next, choose the size and power of the server by looking at the available
flavors::

    $ rack servers flavor list

    ID            Name            RAM    Disk    Swap    VCPUs    RxTxFactor
    2            512MB Standard Instance    512    20    512    1    80
    3            1GB Standard Instance    1024    40    1024    1    120
    4            2GB Standard Instance    2048    80    2048    2    240
    5            4GB Standard Instance    4096    160    2048    2    400
    6            8GB Standard Instance    8192    320    2048    4    600
    7            15GB Standard Instance    15360    620    2048    6    800
    8            30GB Standard Instance    30720    1200    2048    8    1200
    compute1-15        15 GB Compute v1    15360    0    0    8    1250
    compute1-30        30 GB Compute v1    30720    0    0    16    2500
    compute1-4        3.75 GB Compute v1    3840    0    0    2    312.5
    compute1-60        60 GB Compute v1    61440    0    0    32    5000
    compute1-8        7.5 GB Compute v1    7680    0    0    4    625
    general1-1        1 GB General Purpose v1    1024    20    0    1    200
    general1-2        2 GB General Purpose v1    2048    40    0    2    400
    general1-4        4 GB General Purpose v1    4096    80    0    4    800
    general1-8        8 GB General Purpose v1    8192    160    0    8    1600
    io1-120            120 GB I/O v1        122880    40    0    32    10000
    io1-15            15 GB I/O v1        15360    40    0    4    1250
    io1-30            30 GB I/O v1        30720    40    0    8    2500
    io1-60            60 GB I/O v1        61440    40    0    16    5000
    io1-90            90 GB I/O v1        92160    40    0    24    7500
    memory1-120        120 GB Memory v1    122880    0    0    16    5000
    memory1-15        15 GB Memory v1        15360    0    0    2    625
    memory1-240        240 GB Memory v1    245760    0    0    32    10000
    memory1-30        30 GB Memory v1        30720    0    0    4    1250
    memory1-60        60 GB Memory v1        61440    0    0    8    2500
    performance1-1        1 GB Performance    1024    20    0    1    200
    performance1-2        2 GB Performance    2048    40    0    2    400
    performance1-4        4 GB Performance    4096    40    0    4    800
    performance1-8        8 GB Performance    8192    40    0    8    1600
    performance2-120    120 GB Performance    122880    40    0    32    10000
    performance2-15        15 GB Performance    15360    40    0    4    1250
    performance2-30        30 GB Performance    30720    40    0    8    2500
    performance2-60        60 GB Performance    61440    40    0    16    5000
    performance2-90        90 GB Performance    92160    40    0    24    7500

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

    ID        ab95d1d6-27d1-42bb-8cdc-800efcb5fc1f
    AdminPass    k6yfUDkgQfEr

Now you can view the server to make sure the Status is ACTIVE::

   $ rack servers instance list

    ID                    Name        Status    PublicIPv4    PrivateIPv4    Image            Flavor
ab95d1d6-27d1-42bb-8cdc-800efcb5fc1f    devserver    ACTIVE    23.253.50.105    10.209.137.64    09de0a66-3156-48b4-90a5-1cf25a905207    general1-4

To connect to the server with SSH using your public key, use this command::

    $ ssh -i ~/.ssh/id_rsa.pub root@23.253.50.104

