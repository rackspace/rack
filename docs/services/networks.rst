.. _networks:

========
Networks
========

This section lists all cloud networks commands supported by Rackspace CLI.

Commands
--------

All cloud networks commands are based on this syntax::

   rack networks <resource> <action> [command flags]


``network``
~~~~~~~~~~~

All cloud networks network commands use this syntax::

    rack networks network <action> [optional flag]

``list``
^^^^^^^^
Retrieves a list of networks. Optional flags can be used to refine
your search::

    rack networks network list [optional flags]

``create``
^^^^^^^^^^
Creates a network::

    rack networks network create --name <networkName> [optional flags]
    (echo networkName1 && echo networkName2) | rack networks network create --stdin name [optional flags]

``get``
^^^^^^^
Retrieves details on a specified network::

    rack networks network get --id <networkID> [optional flags]
    rack networks network get --name <networkName> [optional flags]
    (echo networkID1 && echo networkID2) | rack networks network get --stdin id [optional flags]

``update``
^^^^^^^^^^
Updates a network::

    rack networks network update --id <networkID> [optional flags]
    rack networks network update --name <networkName> [optional flags]
    (echo networkID1 && echo networkID2) | rack networks network update --stdin id [optional flags]

``delete``
^^^^^^^^^^
Permanently removes a network::

    rack networks network delete --id <networkID> [optional flags]
    rack networks network delete --name <networkName> [optional flags]
    (echo networkID1 && echo networkID2) | rack networks network delete --stdin id [optional flags]

``subnet``
~~~~~~~~~~

Cloud networks subnet commands use this syntax::

    rack networks subnet <action> [optional flags]

``list``
^^^^^^^^
Retrieves a list of subnets. Optional flags can be used to refine your search::

    rack networks subnet list [optional flags]

``create``
^^^^^^^^^^
Creates a subnet::

    rack networks subnet create --network-id <networkID> --cidr <CIDR> --ip-version <IPVersion> [optional flags]

``get``
^^^^^^^
Retrieves details on a specified subnet::

    rack networks subnet get --id <subnetID> [optional flags]
    rack networks subnet get --name <subnetName>] [optional flags]
    (echo subnetID1 && echo subnetID2) | rack networks subnet get --stdin id [optional flags]

``update``
^^^^^^^^^^
Updates a subnet::

    rack networks subnet update --id <subnetID> [optional flags]
    rack networks subnet update --name <subnetName>] [optional flags]

``delete``
^^^^^^^^^^
Permanently deletes a subnet::

    rack networks subnet delete --id <subnetID> [optional flags]
    rack networks subnet delete --name <subnetName>] [optional flags]
    (echo subnetID1 && echo subnetID2) | rack networks subnet delete --stdin id [optional flags]

``port``
~~~~~~~~

Cloud networks port commands use this syntax::

    rack networks port <action> [optional flags]

``list``
^^^^^^^^
Retrieves a list of ports. Optional flags can be used to refine your search::

    rack networks port list [optional flags]

``create``
^^^^^^^^^^
Creates a port::

    rack networks port create --network-id <networkID> [optional flags]

``get``
^^^^^^^
Retrieves details on a specified port::

    rack networks port get --id <portID> [optional flags]
    rack networks port get --name <portName>] [optional flags]
    (echo portID1 && echo portID2) | rack networks port get --stdin id [optional flags]

``update``
^^^^^^^^^^
Updates the information on a port::

    rack networks port update --id <portID> [optional flags]
    rack networks port update --name <portName>] [optional flags]

``delete``
^^^^^^^^^^
Permanently removes a port::

    rack networks port delete --id <portID> [optional flags]
    rack networks port delete --name <portName>] [optional flags]
    (echo portID1 && echo portID2) | rack networks port delete --stdin id [optional flags]

``security-group``
~~~~~~~~~~~~~~~~~~

Cloud network security-group commands use this syntax::

    rack networks security-group <action> [optional flags]

.. note::

    The security groups feature is currently in limited availability. It is available
    only to Managed Infrastructure customers. To use this feature, contact Rackspace S
    Support.

``list``
^^^^^^^^
Retrieves a list of security groups::

    rack networks security-group list [optional flags]

``create``
^^^^^^^^^^
Creates a security group::

    rack networks security-group create --name <securityGroupName> [optional flags]

``get``
^^^^^^^
Retrieves details on a specified security group, including any security group rules::

    rack networks security-group get --id <securityGroupID> [optional flags]
    rack networks security-group get --name <securityGroupName> [optional flags]
    (echo securityGroupID1 && echo securityGroupID2) | rack networks security-group get --stdin id [optional flags]

``delete``
^^^^^^^^^^
Permanently removes a security group and all rules within that security group::

    rack networks security-group delete --id <securityGroupID> [optional flags]
    rack networks security-group delete --name <securityGroupName> [optional flags]
    (echo securityGroupID1 && echo securityGroupID2) | rack networks security-group delete --stdin id [optional flags]

``security-group-rule``
~~~~~~~~~~~~~~~~~~~~~~~

Cloud networks security group rule commands use this syntax::

    rack networks security-group-rule <action> [optional flags]

``list``
^^^^^^^^
Retrieves a list of security group rules with the rules' unique ID::

    rack networks security-group-rule list [optional flags]

``create``
^^^^^^^^^^
Creates a security group rule within a specified security group::

    rack security-group-rule create --security-group-id <securityGroupID> --direction <ingress|egress> --ether-type <ipv4|ipv6> [optional flags]

``get``
^^^^^^^
Retrieves details on a specified security group rule::

    rack networks security-group-rule get --id <securityGroupRuleID> [optional flags]
    (echo securityGroupRuleID1 && echo securityGroupRuleID2) | rack networks security-group-rule get --stdin id [optional flags]

``delete``
^^^^^^^^^^
Permanently deletes a security group rule::

    rack networks security-group-rule delete --id <securityGroupRuleID> [optional flags]
    (echo securityGroupRuleID1 && echo securityGroupRuleID2) | rack networks security-group-rule delete --stdin id [optional flags]
