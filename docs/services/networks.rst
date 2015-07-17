.. _networks:

networks
========

Commands for Rackspace Cloud Networks.

Synopsis
--------

::

   rack networks <resource> <action> [command flags]

Commands
--------

``network``
~~~~~~~~~~~

  Networks Network operations

``list``
^^^^^^^^
Usage::

    rack networks network list [optional flags]

``create``
^^^^^^^^^^
Usage::

    rack networks network create --name <networkName> [optional flags]
    (echo networkName1 && echo networkName2) | rack networks network create --stdin name [optional flags]

``get``
^^^^^^^
Usage::

    rack networks network get --id <networkID> [optional flags]
    rack networks network get --name <networkName> [optional flags]
    (echo networkID1 && echo networkID2) | rack networks network get --stdin id [optional flags]

``update``
^^^^^^^^^^
Usage::

    rack networks network update --id <networkID> [optional flags]
    rack networks network update --name <networkName> [optional flags]
    (echo networkID1 && echo networkID2) | rack networks network update --stdin id [optional flags]

``delete``
^^^^^^^^^^
Usage::

    rack networks network delete --id <networkID> [optional flags]
    rack networks network delete --name <networkName> [optional flags]
    (echo networkID1 && echo networkID2) | rack networks network delete --stdin id [optional flags]

``subnet``
~~~~~~~~~~

  Networks Subnet operations

``list``
^^^^^^^^
Usage::

    rack networks subnet list [optional flags]

``create``
^^^^^^^^^^
Usage::
    rack networks subnet create --network-id <networkID> --cidr <CIDR> --ip-version <IPVersion> [optional flags]

``get``
^^^^^^^
Usage::

    rack networks subnet get --id <subnetID> [optional flags]
    rack networks subnet get --name <subnetName>] [optional flags]
    (echo subnetID1 && echo subnetID2) | rack networks subnet get --stdin id [optional flags]

``update``
^^^^^^^^^^
Usage::

    rack networks subnet update --id <subnetID> [optional flags]
    rack networks subnet update --name <subnetName>] [optional flags]

``delete``
^^^^^^^^^^
Usage::

    rack networks subnet delete --id <subnetID> [optional flags]
    rack networks subnet delete --name <subnetName>] [optional flags]
    (echo subnetID1 && echo subnetID2) | rack networks subnet delete --stdin id [optional flags]

``port``
~~~~~~~~

  Networks Port operations

``list``
^^^^^^^^
Usage::

    rack networks port list [optional flags]

``create``
^^^^^^^^^^
Usage::

    rack networks port create --network-id <networkID> [optional flags]

``get``
^^^^^^^
Usage::

    rack networks port get --id <portID> [optional flags]
    rack networks port get --name <portName>] [optional flags]
    (echo portID1 && echo portID2) | rack networks port get --stdin id [optional flags]

``update``
^^^^^^^^^^
Usage::

    rack networks port update --id <portID> [optional flags]
    rack networks port update --name <portName>] [optional flags]

``delete``
^^^^^^^^^^
Usage::

    rack networks port delete --id <portID> [optional flags]
    rack networks port delete --name <portName>] [optional flags]
    (echo portID1 && echo portID2) | rack networks port delete --stdin id [optional flags]

``security-group``
~~~~~~~~~~~~~~~~~~

  Networks Security Group operations

``list``
^^^^^^^^
Usage::

    rack networks security-group list [flags]

``create``
^^^^^^^^^^
Usage::

    rack networks security-group create --name <securityGroupName> [optional flags]

``get``
^^^^^^^
Usage::

    rack networks security-group get --id <securityGroupID> [optional flags]
    rack networks security-group get --name <securityGroupName> [optional flags]
    (echo securityGroupID1 && echo securityGroupID2) | rack networks security-group get --stdin id [optional flags]

``delete``
^^^^^^^^^^
Usage::

    rack networks security-group delete --id <securityGroupID> [optional flags]
    rack networks security-group delete --name <securityGroupName> [optional flags]
    (echo securityGroupID1 && echo securityGroupID2) | rack networks security-group delete --stdin id [optional flags]

``security-group-rule``
~~~~~~~~~~~~~~~~~~~~~~~

  Networks Security Group Rule operations

``list``
^^^^^^^^
Usage::

    rack networks security-group-rule list [optional flags]

``create``
^^^^^^^^^^
Usage::

    rack security-group-rule create --security-group-id <securityGroupID> --direction <ingress|egress> --ether-type <ipv4|ipv6> [optional flags]

``get``
^^^^^^^
Usage::

    rack networks security-group-rule get --id <securityGroupRuleID> [optional flags]
    (echo securityGroupRuleID1 && echo securityGroupRuleID2) | rack networks security-group-rule get --stdin id [optional flags]

``delete``
^^^^^^^^^^
Usage::

    rack networks security-group-rule delete --id <securityGroupRuleID> [optional flags]
    (echo securityGroupRuleID1 && echo securityGroupRuleID2) | rack networks security-group-rule delete --stdin id [optional flags]
