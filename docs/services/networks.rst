.. _networks:

=================
Networks commands
=================

This section lists the commands for interacting with Rackspace Cloud Networks.

All ``networks`` commands are based on the following syntax::

   rack networks <subcommand> <action> [command flags]

*Command flags* enable you to customize certain attributes of the command, such as using ``--name`` to name a network. To display a list of command flags specific to the command, type ``rack networks <subcommand> <action> --help``.

The following sections describe the ``networks`` subcommands and the actions associated with them.

Network
-------

The ``network`` subcommand provides information about and performs actions on the networks in Cloud Networks. The ``network`` subcommand uses the following syntax::

    rack networks network <action> [optional flag]

The following sections describe the actions that you can perform on the ``network`` subcommand and provide example responses.

``list``
~~~~~~~~
Retrieves a list of networks.

::

    rack networks network list [optional flags]

**Response**

.. code::

    $ rack networks network list
    ID					                          Name				                    Up	  Status	Shared	Tenant ID
    6843cc43-7dc5-4164-96c0-b7f404fd6120	My Private Network		          true	ACTIVE	false	  661145
    bb3816e1-514d-4543-bfda-358e81a2c8a5	test				                    true	ACTIVE	false	  661145
    e6fba9eb-5211-4637-bf1a-cdb4c04a2845	My Test Network			            true	ACTIVE	false	  661145
    ec962e06-b566-4911-9b9f-f9a45796921c	My Private Network for Vy Class	true	ACTIVE	false	  661145

``create``
~~~~~~~~~~
Creates a network with the specified name.

::

    rack networks network create --name <networkName> [optional flags]
    (echo networkName1 && echo networkName2) | rack networks network create --stdin name [optional flags]

**Response**

.. code::

    $ rack networks network create --name RackCLItest
    ID		4d83cf66-26be-44cc-b344-034e8e58c769
    Name		RackCLItest
    Up		true
    Status		ACTIVE
    Shared		false
    TenantID	661145

``get``
~~~~~~~
Retrieves details about a network, which you can specify by ID or name. 

::

    rack networks network get --id <networkID> [optional flags]
    rack networks network get --name <networkName> [optional flags]
    (echo networkID1 && echo networkID2) | rack networks network get --stdin id [optional flags]

**Response**

.. code::

    $ rack networks network get --name RackCLItest
    ID		4d83cf66-26be-44cc-b344-034e8e58c769
    Name		RackCLItest
    Up		true
    Status		ACTIVE
    Subnets
    Shared		false
    TenantID	661145

``update``
~~~~~~~~~~
Updates a network, which you can specify by ID or name. 

::

    rack networks network update --id <networkID> [optional flags]
    rack networks network update --name <networkName> [optional flags]
    (echo networkID1 && echo networkID2) | rack networks network update --stdin id [optional flags]

**Response**

.. code::

    $ rack networks network update --id 4d83cf66-26be-44cc-b344-034e8e58c769 --up false
    ID		4d83cf66-26be-44cc-b344-034e8e58c769
    Name		RackCLItest
    Up		false
    Status		ACTIVE
    Shared		false
    TenantID	661145

``delete``
~~~~~~~~~~
Permanently deletes a network, which you can specify by ID or name. 

::

    rack networks network delete --id <networkID> [optional flags]
    rack networks network delete --name <networkName> [optional flags]
    (echo networkID1 && echo networkID2) | rack networks network delete --stdin id [optional flags]

**Response**

.. code::

    $ rack networks network delete --name RackCLItest
    Successfully deleted network [4d83cf66-26be-44cc-b344-034e8e58c769]

Subnet
------

The ``subnet`` subcommand provides information about and performs actions on the subnets in Cloud Networks. The ``subnet`` subcommand uses the following syntax::

    rack networks subnet <action> [optional flags]

The following sections describe the actions that you can perform on the ``subnet`` subcommand and provide example responses.

``list``
~~~~~~~~
Retrieves a list of subnets. 

::

    rack networks subnet list [optional flags]

``create``
~~~~~~~~~~
Creates a subnet with the specified details.

::

    rack networks subnet create --network-id <networkID> --cidr <CIDR> --ip-version <4|6> [optional flags]

**Response**

.. code::

    $ rack networks subnet create --network-id 8694604a-eb48-4c69-8fce-ce8fa647fa30 --cidr 192.0.2.0/24 --ip-version 4
    ID		ed3161fa-f1d2-45e5-bd4d-4c5340ad0027
    Name
    NetworkID	8694604a-eb48-4c69-8fce-ce8fa647fa30
    CIDR		192.0.2.0/24
    EnableDHCP	false
    GatewayIP
    DNSNameservers
    AllocationPools0Start   192.0.2.1
    AllocationPools0End   192.0.2.254

``get``
~~~~~~~
Retrieves details about a subnet, which you can specify by ID or name. 

::

    rack networks subnet get --id <subnetID> [optional flags]
    rack networks subnet get --name <subnetName>] [optional flags]
    (echo subnetID1 && echo subnetID2) | rack networks subnet get --stdin id [optional flags]

**Response**

.. code::

    $ rack networks subnet get --id ed3161fa-f1d2-45e5-bd4d-4c5340ad0027
    ID		ed3161fa-f1d2-45e5-bd4d-4c5340ad0027
    Name
    NetworkID	8694604a-eb48-4c69-8fce-ce8fa647fa30
    CIDR		192.0.2.0/24
    EnableDHCP	false
    GatewayIP
    DNSNameservers
    AllocationPools0Start   192.0.2.1
    AllocationPools0End   192.0.2.254
    HostRoutes	[]

``update``
~~~~~~~~~~
Updates a subnet, which you can specify by ID or name. 

::

    rack networks subnet update --id <subnetID> [optional flags]
    rack networks subnet update --name <subnetName>] [optional flags]

**Response**

.. code::

    $ rack networks subnet update --id ed3161fa-f1d2-45e5-bd4d-4c5340ad0027 --rename CLIsub
    ID		ed3161fa-f1d2-45e5-bd4d-4c5340ad0027
    Name		CLIsub
    NetworkID	8694604a-eb48-4c69-8fce-ce8fa647fa30
    CIDR		192.0.2.0/24
    EnableDHCP	false
    GatewayIP
    DNSNameservers
    AllocationPools0Start   192.0.2.1
    AllocationPools0End   192.0.2.254

``delete``
~~~~~~~~~~
Permanently deletes a subnet, which you can specify by ID or name. 

::

    rack networks subnet delete --id <subnetID> [optional flags]
    rack networks subnet delete --name <subnetName>] [optional flags]
    (echo subnetID1 && echo subnetID2) | rack networks subnet delete --stdin id [optional flags]

**Response**

.. code::

    $ rack networks subnet delete --name CLIsub
    Successfully deleted subnet [ed3161fa-f1d2-45e5-bd4d-4c5340ad0027]

Port
----

The ``port`` subcommand provides information about and performs actions on the ports in Cloud Networks. The ``port`` subcommand uses the following syntax::

    rack networks port <action> [optional flags]

The following sections describe the actions that you can perform on the ``port`` subcommand and provide example responses.

``list``
~~~~~~~~
Retrieves a list of ports.

::

    rack networks port list [optional flags]

**Response**

.. code::

    $ rack networks port list --network-id e6fba9eb-5211-4637-bf1a-cdb4c04a2845
    ID					                        Name	NetworkID				                      Status	MACAddress		DeviceID
    3456c2b0-7bee-40b4-ad0b-b4f3385fb36c		  e6fba9eb-5211-4637-bf1a-cdb4c04a2845	ACTIVE	BC:76:4E:05:FF:1B

``create``
~~~~~~~~~~
Creates a port that is associated with the specified network ID. 

::

    rack networks port create --network-id <networkID> [optional flags]

**Response**

.. code::

    $ rack networks port create --network-id e6fba9eb-5211-4637-bf1a-cdb4c04a2845
    ID		3456c2b0-7bee-40b4-ad0b-b4f3385fb36c
    Name
    NetworkID	e6fba9eb-5211-4637-bf1a-cdb4c04a2845
    Status		ACTIVE
    MACAddress	BC:76:4E:05:FF:1B
    DeviceID
    DeviceOwner
    Up		false
    FixedIPs0:SubnetID    e3cdb6bd-f308-4c15-93db-7638dd995111
    FixedIPs0:IPAddress   192.168.4.3
    SecurityGroups
    TenantID	661145

``get``
~~~~~~~
Retrieves details about a port, which you can specify by ID or name. 

::

    rack networks port get --id <portID> [optional flags]
    rack networks port get --name <portName>] [optional flags]
    (echo portID1 && echo portID2) | rack networks port get --stdin id [optional flags]

**Response**

.. code::

    $ rack networks port get --id 3456c2b0-7bee-40b4-ad0b-b4f3385fb36c
    ID		3456c2b0-7bee-40b4-ad0b-b4f3385fb36c
    Name
    NetworkID	e6fba9eb-5211-4637-bf1a-cdb4c04a2845
    Status		ACTIVE
    MACAddress	BC:76:4E:05:FF:1B
    DeviceID
    DeviceOwner
    Up		false
    FixedIPs0:SubnetID    e3cdb6bd-f308-4c15-93db-7638dd995111
    FixedIPs0:IPAddress   192.168.4.3
    SecurityGroups
    TenantID	661145

``update``
~~~~~~~~~~
Updates the information about a port, which you can specify by ID or name. 

::

    rack networks port update --id <portID> [optional flags]
    rack networks port update --name <portName>] [optional flags]

**Response**

.. code::

    $ rack networks port update --id 3456c2b0-7bee-40b4-ad0b-b4f3385fb36c --rename CLIport
    ID		3456c2b0-7bee-40b4-ad0b-b4f3385fb36c
    Name	CLIport
    NetworkID	e6fba9eb-5211-4637-bf1a-cdb4c04a2845
    Status		ACTIVE
    MACAddress	BC:76:4E:05:FF:1B
    DeviceID
    DeviceOwner
    Up		false
    FixedIPs0:SubnetID    e3cdb6bd-f308-4c15-93db-7638dd995111
    FixedIPs0:IPAddress   192.168.4.3
    SecurityGroups
    TenantID	661145

``delete``
~~~~~~~~~~
Permanently deletes a port, which you can specify by ID or name. 

::

    rack networks port delete --id <portID> [optional flags]
    rack networks port delete --name <portName>] [optional flags]
    (echo portID1 && echo portID2) | rack networks port delete --stdin id [optional flags]

**Response**

.. code::

    $ rack networks port delete --name CLIport
    Successfully deleted port [3456c2b0-7bee-40b4-ad0b-b4f3385fb36c]

Security-group
--------------

The ``security-group`` subcommand provides information about and performs actions on the security groups in Cloud Networks. The ``security-group`` subcommand uses the following syntax::

    rack networks security-group <action> [optional flags]

The following sections describe the actions that you can perform on the ``security-group`` subcommand and provide example responses.

.. note::

    The security groups feature is currently in limited availability. It is available only to Managed Infrastructure customers. To use this feature, contact Rackspace Support.

``list``
~~~~~~~~
Retrieves a list of security groups.

::

    rack networks security-group list [optional flags]

**Response**

.. code::

    $ rack networks security-group list
    ID					                          Name	  TenantID
    928fb119-9c69-4f9f-8da5-8387fd923863	CLIsec	661145

``create``
~~~~~~~~~~
Creates a security group with the specified name.

::

    rack networks security-group create --name <securityGroupName> [optional flags]

**Response**

.. code::

    $ rack networks security-group create --name CLIsec
    ID	928fb119-9c69-4f9f-8da5-8387fd923863
    Name	CLIsec

``get``
~~~~~~~
Retrieves details about a security group, including any security group rules. You can specify the security group by ID or name.

::

    rack networks security-group get --id <securityGroupID> [optional flags]
    rack networks security-group get --name <securityGroupName> [optional flags]
    (echo securityGroupID1 && echo securityGroupID2) | rack networks security-group get --stdin id [optional flags]

**Response**

.. code::

    $ rack networks security-group get --name CLIsec
    ID	928fb119-9c69-4f9f-8da5-8387fd923863
    Name	CLIsec
    TenantID661145
    Rules0:EtherTypeIPv4
    Rules0:Protocol
    Rules0:ID	ff0029e9-f09d-4ddd-889f-36f9c2ff316b
    Rules0:Directioningress


``delete``
~~~~~~~~~~
Permanently deletes a security group and all rules within that security group. You can specify the security group by ID or name.

::

    rack networks security-group delete --id <securityGroupID> [optional flags]
    rack networks security-group delete --name <securityGroupName> [optional flags]
    (echo securityGroupID1 && echo securityGroupID2) | rack networks security-group delete --stdin id [optional flags]

**Response**

.. code::

    $ rack networks security-group delete --name CLIsec
    Successfully deleted security group [928fb119-9c69-4f9f-8da5-8387fd923863]


Security-group-rule
-------------------

The ``security-group-rule`` subcommand provides information about and performs actions on the security group rules in Cloud Networks. The ``security-group-rule`` subcommand uses the following syntax::

    rack networks security-group-rule <action> [optional flags]

The following sections describe the actions that you can perform on the ``security-group-rule`` subcommand and provide example responses.

``list``
~~~~~~~~
Retrieves a list of security group rules.

::

    rack networks security-group-rule list [optional flags]

**Response**

.. code::

    $ rack networks security-group-rule list
    ID					                          Direction	EtherType	PortRangeMin	PortRangeMax	Protocol	SecurityGroupID
    a84602ac-8f79-4fe5-9329-2ceebaa958da	ingress		IPv4		  80		        80		        TCP		    928fb119-9c69-4f9f-8da5-8387fd923863

``create``
~~~~~~~~~~
Creates a security group rule within a specified security group.

::

    rack networks security-group-rule create --security-group-id <securityGroupID> --direction <ingress|egress> --ether-type <ipv4|ipv6> [optional flags]

**Response**

.. code::

    $ rack networks security-group-rule create --security-group-id 928fb119-9c69-4f9f-8da5-8387fd923863 --direction ingress  --ether-type ipv4 --port-range-min 80 --port-range-max 80 --protocol tcp
    ID		a84602ac-8f79-4fe5-9329-2ceebaa958da
    Direction	ingress
    EtherType	IPv4
    PortRangeMin	80
    PortRangeMax	80
    Protocol	TCP
    SecurityGroupID	928fb119-9c69-4f9f-8da5-8387fd923863

``get``
~~~~~~~
Retrieves details about the specified security group rule.

::

    rack networks security-group-rule get --id <securityGroupRuleID> [optional flags]
    (echo securityGroupRuleID1 && echo securityGroupRuleID2) | rack networks security-group-rule get --stdin id [optional flags]

**Response**

.. code::

    $ rack networks security-group-rule get --id a84602ac-8f79-4fe5-9329-2ceebaa958da
    ID		a84602ac-8f79-4fe5-9329-2ceebaa958da
    Direction	ingress
    EtherType	IPv4
    PortRangeMin	80
    PortRangeMax	80
    Protocol	TCP
    SecurityGroupID	928fb119-9c69-4f9f-8da5-8387fd923863
    RemoteGroupID
    RemoteIPPrefix
    TenantID	661145

``delete``
~~~~~~~~~~
Permanently deletes a security group rule.

::

    rack networks security-group-rule delete --id <securityGroupRuleID> [optional flags]
    (echo securityGroupRuleID1 && echo securityGroupRuleID2) | rack networks security-group-rule delete --stdin id [optional flags]

**Response**

.. code::

    $ rack networks security-group-rule delete --id a84602ac-8f79-4fe5-9329-2ceebaa958da
    Successfully deleted security group rule [a84602ac-8f79-4fe5-9329-2ceebaa958da]
