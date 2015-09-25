.. _orchestration:

=======
Orchestration
=======

This section lists the commands for interacting with Cloud Orchestration.

Commands
--------

All ``orchestration`` commands are based on this syntax::

   rack orchestration <subservice> <action> [command flags]

*Command flags* allow you to customize certain attributes of the command,
such as ``--name`` to name a stack. Type ``rack orchestration <subservice> <action> --help``
to bring up a list *command flags* specific to the command.

**BuildInfo**
~~~~~~~~~~~~

Retrieves the build information of the orchestration service::

    rack orchestration build-info get [optional flags]
**Response**

.. code::

    $ rack orchestration build-info get
    API		2015.l3-20150903-1517
    Engine		2015.l3-20150903-1517
    FusionAPI	l1-20150622-17c7bae-141

**Stack**
~~~~~~~~~~~~

Stack commands use this syntax::

    rack orchestration stack <action> [command flags]

``list``
^^^^^^^^
Retrieves a list of stacks::

    rack orchestration stack list [optional flags]

**Response**

.. code::

    $ rack orchestration stack list
    ID					                    Name					            Status		        CreationTime
    82f5e08a-6429-4687-a36d-e25733d04a26	vijendar-teststack			        CREATE_COMPLETE	    2015-09-11 16:32:33 +0000 UTC
    b7eceee2-53af-44d4-8265-ce7209f081a4	drago_heat_dev3_DO_NOT_DELETE		CREATE_COMPLETE	    2015-08-31 06:38:13 +0000 UTC
    373f8758-0777-414e-b963-d0cd8f16c019	pm-test_DO_NOT_DELETE			    CREATE_COMPLETE	    2015-08-26 18:27:59 +0000 UTC
    54302da9-a2b7-4378-af93-81763bf8aa4f	drago_heat_dev_env3_DO_NOT_DELETE	CREATE_COMPLETE	    2015-08-25 18:01:56 +0000 UTC
    0a8df0f3-8a58-4634-9654-9eba6629bd52	pm-devstack_DO_NOT_DELETE		    CREATE_FAILED	    2015-08-21 06:57:33 +0000 UTC
    5e15d018-374c-48fa-b5b2-96e8d3bb41d3	pm-test-heatdev-DO_NOT_DELETE		CREATE_COMPLETE     2015-08-17 13:58:39 +0000 UTC
    5b56395a-4e8b-4389-bd44-a123030c7c9c	pm_test_bug_DO_NOT_DELETE		    DELETE_FAILED	    2015-07-24 14:38:00 +0000 UTC

``create``
^^^^^^^^^^
Creates a stack::

    rack orchestration stack create --name <stackName> [optional flags]
    (echo stackName1 && echo stackName2) | rack orchestration stack create --stdin name [optional flags]

In order for this command to work, you must select a template using either the
``--template-file`` or ``--template-url`` flags with the file or URL of the template you wish to use.

**Response**

.. code::

    $ rack orchestration stack create --name RackTest --template-file mytemplate.yaml
    ID					                    Name					            Status		        CreationTime
    97c2e5a2-7d8c-4c7e-8fcc-eac122634f34    RackTest                            CREATE_IN_PROGRESS  2015-09-11 22:01:07 +0000 UTC
    5b56395a-4e8b-4389-bd44-a123030c7c9c	pm_test_bug_DO_NOT_DELETE		    DELETE_FAILED	    2015-07-24 14:38:00 +0000 UTC


.. note::

    The details of the stack you created can be found by using the ``get``
    commmand described below.

``get``
^^^^^^^
Retrieves details of a specified stack::

    rack orchestration stack get --id <stackID> [optional flags]
    rack orchestration stack get --name <stackName> [optional flags]
    (echo stackName1 && echo stackName2) | rack orchestration stack get --stdin name [optional flags]

**Response**

.. code::

    $ rack orchestration stack get --name RackTest
    Capabilities		      []
    CreationTime		      2015-09-11 22:01:07 +0000 UTC
    Description		          No description
    DisableRollback		      true
    ID			              97c2e5a2-7d8c-4c7e-8fcc-eac122634f34
    Links			          [{
    			                "Href": "https://iad.orchestration.api.rackspacecloud.com/v1/897686/stacks/RackTest/97c2e5a2-7d8c-4c7e-8fcc-eac122634f34",
    			                "Rel": "self"
    			              }]
    NotificiationTopics	      <nil>
    Outputs			          []
    Parameters		          {
    			                "OS::project_id": "897686",
    			                "OS::stack_id": "97c2e5a2-7d8c-4c7e-8fcc-eac122634f34",
    			                "OS::stack_name": "RackTest",
    			                "flavor": "******"
    			              }
    Name			          RackTest
    Status			          CREATE_COMPLETE
    StatusReason		      Stack CREATE completed successfully
    TemplateDescription	      No description
    Timeout			          None
    Tags			          []
    UpdatedTime		          None

``update``
^^^^^^^^^^
Updates the stack using a provided template::

    rack orchestration stack update --id <stackID> [optional flags]
    rack orchestration stack update --name <stackName> [optional flags]


**Response**

.. code::

    $ rack orchestration stack update --name RackTest --template-file myNewTemplate.yaml
    ID					                    Name					            Status		        CreationTime
    97c2e5a2-7d8c-4c7e-8fcc-eac122634f34	RackTest				            UPDATE_IN_PROGRESS	2015-09-11 22:01:07 +0000 UTC
    5b56395a-4e8b-4389-bd44-a123030c7c9c	pm_test_bug_DO_NOT_DELETE		    DELETE_FAILED	    2015-07-24 14:38:00 +0000 UTC


``delete``
^^^^^^^^^^
Deletes a stack::

    rack orchestration stack delete --id <stackID> [optional flags]
    rack orchestration stack delete --name <stackName> [optional flags]
    (echo stackName1 && echo stackName2) | rack orchestration stack delete --stdin name [optional flags]

**Response**

.. code::

    $ rack orchestration stack delete --name RackTest
    ID					                    Name					            Status		        CreationTime
    97c2e5a2-7d8c-4c7e-8fcc-eac122634f34	RackTest				            DELETE_IN_PROGRESS	2015-09-11 22:01:07 +0000 UTC
    5b56395a-4e8b-4389-bd44-a123030c7c9c	pm_test_bug_DO_NOT_DELETE		    DELETE_FAILED	    2015-07-24 14:38:00 +0000 UTC


``preview``
^^^^^^^^^^
Preview shows the number and type of resources that will be created by a template::

    rack  orchestration stack preview --name <stackName> [--template-file <templateFile> | --template-url <templateURL>] [optional flags]
    (echo stackName1 && echo stackName2) | rack  orchestration stack preview --stdin name [--template-file <templateFile> | --template-url <templateURL>] [optional flags]

**Response**

.. code::

    $ rack orchestration stack preview --template-file my_nova.yaml --name RackTest
    Capabilities		[]
    CreationTime		2015-09-11 22:55:51.60336 +0000 UTC
    Description		No description
    DisableRollback		true
    ID			None
    Links			[{
    			  "Href": "https://iad.orchestration.api.rackspacecloud.com/v1/897686/stacks/RackTest/None",
    			  "Rel": "self"
    			}]
    NotificiationTopics	<nil>
    Parameters		{
    			  "OS::project_id": "897686",
    			  "OS::stack_id": "None",
    			  "OS::stack_name": "RackTest",
    			  "flavor": "******"
    			}
    Resources		[
    			  {
    			    "attributes": {
    			      "accessIPv4": "",
    			      "accessIPv6": "",
    			      "addresses": "",
    			      "console_urls": "",
    			      "first_address": "",
    			      "instance_name": "",
    			      "name": "test-server-2",
    			      "networks": ""
    			    },
    			    "creation_time": null,
    			    "description": "",
    			    "metadata": {},
    			    "physical_resource_id": "",
    			    "properties": {
    			      "admin_pass": null,
    			      "admin_user": null,
    			      "availability_zone": null,
    			      "block_device_mapping": null,
    			      "block_device_mapping_v2": null,
    			      "config_drive": null,
    			      "diskConfig": null,
    			      "flavor": "2 GB General Purpose v1",
    			      "flavor_update_policy": "RESIZE",
    			      "image": "Debian 7 (Wheezy) (PVHVM)",
    			      "image_update_policy": "REBUILD",
    			      "key_name": null,
    			      "metadata": null,
    			      "name": "test-server-2",
    			      "networks": null,
    			      "personality": {},
    			      "reservation_id": null,
    			      "scheduler_hints": null,
    			      "security_groups": [],
    			      "software_config_transport": "POLL_TEMP_URL",
    			      "user_data": "",
    			      "user_data_format": "RAW"
    			    },
    			    "required_by": [],
    			    "resource_action": "INIT",
    			    "resource_identity": {
    			      "path": "/resources/test_server",
    			      "stack_id": "None",
    			      "stack_name": "RackTest",
    			      "tenant": "897686"
    			    },
    			    "resource_name": "test_server",
    			    "resource_status": "COMPLETE",
    			    "resource_status_reason": "",
    			    "resource_type": "OS::Nova::Server",
    			    "stack_identity": {
    			      "path": "",
    			      "stack_id": "None",
    			      "stack_name": "RackTest",
    			      "tenant": "897686"
    			    },
    			    "stack_name": "RackTest",
    			    "updated_time": null
    			  }
    			]
    Name			RackTest
    TemplateDescription	No description
    Timeout			None
    UpdatedTime		None

``abandon``
^^^^^^^^^^^
Abandons the stack. This will delete the record of the stack from orchestration, but
will not delete any of the underlying resources::

    rack orchestration stack abandon --id <stackID> [optional flags]
    rack orchestration stack abandon --name <stackName> [optional flags]

Prints an JSON representation of the stack to stdout or a file on success. This
can be used in the ``adopt`` command to create a new stack with the existing
resources.

**Response**

.. code::

    $ rack orchestration stack abandon --name RackTest
    Status			COMPLETE
    Name			RackTest
    Template		{
    			  "heat_template_version": "2014-10-16",
    			  "parameters": {
    			    "flavor": {
    			      "default": 4353,
    			      "description": "Flavor for the server to be created",
    			      "hidden": true,
    			      "type": "string"
    			    }
    			  },
    			  "resources": {
    			    "test_server": {
    			      "properties": {
    			        "flavor": "2 GB General Purpose v1",
    			        "image": "Debian 7 (Wheezy) (PVHVM)",
    			        "name": "test-server-2"
    			      },
    			      "type": "OS::Nova::Server"
    			    }
    			  }
    			}
    Action			CREATE
    ID			2c4f91a6-228a-40f8-a842-d20bef420ad0
    Resources		{
    			  "test_server": {
    			    "action": "CREATE",
    			    "metadata": {},
    			    "name": "test_server",
    			    "resource_data": {},
    			    "resource_id": "69c99fc6-a856-4e37-ac28-9e19de6cb674",
    			    "status": "COMPLETE",
    			    "type": "OS::Nova::Server"
    			  }
    			}
    Files			{}
    StackUserProjectID	897686
    ProjectID		897686
    Environment		{
    			  "encrypted_param_names": [],
    			  "parameter_defaults": {},
    			  "parameters": {},
    			  "resource_registry": {
    			    "resources": {}
    			  }
    			}

``adopt``
^^^^^^^^^^
Creates a stack without creating any resources; existing resources are used
instead::

    rack orchestration stack adopt --name stackName --adopt-file adoptFile [optional flags]

**Response**

.. code::

    $ rack orchestration stack adopt --name RackTest --adopt-file abandon.yaml
    ID					                    Name					            Status		        CreationTime
    43cedc45-e188-4e49-88a9-728b2126586c	RackTest				            ADOPT_COMPLETE	    2015-09-11 23:40:18 +0000 UTC
    5b56395a-4e8b-4389-bd44-a123030c7c9c	pm_test_bug_DO_NOT_DELETE		    DELETE_FAILED	    2015-07-24 14:38:00 +0000 UTC

**Resource**
~~~~~~~~~~~~

Resource commands use this syntax::

    rack orchestration resource <action> [command flags]

``list``
^^^^^^^^
Retrieves a list of resources for a given stack::

    rack orchestration resource list --name <stackName> [optional flags]
    rack orchestration resource list --id <stackID> [optional flags]
    (echo stackName1 && echo stackName2)  | rack orchestration resource list --stdin name [optional flags]

**Response**

.. code::

    $ rack orchestration resource list --name RackTest
    Name		PhysicalID				                Type			    Status		               UpdatedTime
    test_server	f075a7c1-28ef-4699-9046-383098134902	OS::Nova::Server	CREATE_COMPLETE	2015-09-12 16:37:49 +0000 UTC


``get``
^^^^^^^
Retrieves details of a specified resource in a stack::

    rack orchestration resource get --id <stackID> --resource <resourceName> [optional flags]
    rack orchestration resource get --name <stackName> --resource <resourceName> [optional flags]

**Response**

.. code::

    $ rack orchestration resource get --name RackTest --resource test_server
    Name		test_server
    PhysicalID	f075a7c1-28ef-4699-9046-383098134902
    Type		OS::Nova::Server
    Status		ADOPT_COMPLETE
    UpdatedTime	2015-09-13 04:20:24 +0000 UTC
    Links		[{
    		  "Href": "https://iad.orchestration.api.rackspacecloud.com/v1/897686/stacks/RackTest/52f8681a-4d5b-46ef-b643-6e945ae85d16/resources/test_server",
    		  "Rel": "self"
    		} {
    		  "Href": "https://iad.orchestration.api.rackspacecloud.com/v1/897686/stacks/RackTest/52f8681a-4d5b-46ef-b643-6e945ae85d16",
    		  "Rel": "stack"
    		}]
    Attributes	{
    		  "OS-DCF:diskConfig": "MANUAL",
    		  "OS-EXT-STS:power_state": 1,
    		  "OS-EXT-STS:task_state": null,
    		  "OS-EXT-STS:vm_state": "active",
    		  "RAX-PUBLIC-IP-ZONE-ID:publicIPZoneId": "4eefdfdcc0c65b6459cb32da6041e307c8b3a7ba1f8f68aa2ed5a746",
    		  "accessIPv4": "104.239.165.61",
    		  "accessIPv6": "2001:4802:7805:101:be76:4eff:fe20:ded8",
    		  "addresses": {
    		    "private": [
    		      {
    		        "addr": "10.209.67.179",
    		        "version": 4
    		      }
    		    ],
    		    "public": [
    		      {
    		        "addr": "104.239.165.61",
    		        "version": 4
    		      },
    		      {
    		        "addr": "2001:4802:7805:101:be76:4eff:fe20:ded8",
    		        "version": 6
    		      }
    		    ]
    		  },
    		  "config_drive": "",
    		  "created": "2015-09-12T16:37:51Z",
    		  "flavor": {
    		    "id": "general1-2",
    		    "links": [
    		      {
    		        "href": "https://iad.servers.api.rackspacecloud.com/897686/flavors/general1-2",
    		        "rel": "bookmark"
    		      }
    		    ]
    		  },
    		  "hostId": "7a4a76cfba0997a0a60d4c57f4c1b8da08b65706a4eb7b66762136c6",
    		  "id": "f075a7c1-28ef-4699-9046-383098134902",
    		  "image": {
    		    "id": "c934d497-7b45-4764-ac63-5b67e1458a20",
    		    "links": [
    		      {
    		        "href": "https://iad.servers.api.rackspacecloud.com/897686/images/c934d497-7b45-4764-ac63-5b67e1458a20",
    		        "rel": "bookmark"
    		      }
    		    ]
    		  },
    		  "key_name": null,
    		  "links": [
    		    {
    		      "href": "https://iad.servers.api.rackspacecloud.com/v2/897686/servers/f075a7c1-28ef-4699-9046-383098134902",
    		      "rel": "self"
    		    },
    		    {
    		      "href": "https://iad.servers.api.rackspacecloud.com/897686/servers/f075a7c1-28ef-4699-9046-383098134902",
    		      "rel": "bookmark"
    		    }
    		  ],
    		  "metadata": {
    		    "rax_service_level_automation": "Complete"
    		  },
    		  "name": "test-server-2",
    		  "progress": 100,
    		  "status": "ACTIVE",
    		  "tenant_id": "897686",
    		  "updated": "2015-09-12T16:38:18Z",
    		  "user_id": "5c11b69d82cf4313b7a8b173b799a0ef"
    		}
    CreationTime	2015-09-13 04:20:24 +0000 UTC
    Description
    LogicalID	test_server

``get-schema``
^^^^^^^^^^
Shows the interface schema for a specified resource type::

    rack orchestration resource get-schema --resource <resourceName> [optional flags]
    (echo stackName1 && echo stackName2) | rack orchestration resource get-schema --stdin name [optional flags]

This schema describes the properties that can be set on the resource, their
types, constraints, descriptions, and default values. Additionally, the
resource attributes and their descriptions are provided.

**Response**

.. code::

    $ rack orchestration resource get-schema --type OS::Heat::None
    Attributes	{
		  "show": {
		    "description": "Detailed information about resource.",
		    "type": "map"
		  }
		}
    Properties	{}
    ResourceType	OS::Heat::None
    SupportStatus	{
		  "message": null,
		  "previous_status": null,
		  "status": "SUPPORTED",
		  "version": "5.0.0"
		}

``get-template``
^^^^^^^^^^
Shows a template representation for specified resource type::

    rack orchestration resource get-template --type <resourceType> [optional flags]
    (echo stackName1 && echo stackName2) | rack orchestration resource get-template --stdin type [optional flags]

**Response**

.. code::

    $ rack orchestration resource get-template --type OS::Heat::None
    {
      "Description": "Initial template of NoneResource",
      "HeatTemplateFormatVersion": "2012-12-12",
      "Outputs": {
        "show": {
          "Description": "Detailed information about resource.",
          "Value": "{\"Fn::GetAtt\": [\"NoneResource\", \"show\"]}"
        }
      },
      "Parameters": {},
      "Resources": {
        "NoneResource": {
          "Properties": {},
          "Type": "OS::Heat::None"
        }
      }
    }

``list-types``
^^^^^^^^
Lists the supported template resource types::

    rack orchestration resource list-types [optional flags]

**Response**

.. code::

    $ rack orchestration resource list-types
    ResourceType
    AWS::CloudFormation::WaitCondition
    AWS::CloudFormation::WaitConditionHandle
    AWS::EC2::Instance
    AWS::ElasticLoadBalancing::LoadBalancer
    DockerInc::Docker::Container
    OS::Cinder::Volume
    OS::Cinder::VolumeAttachment
    OS::Heat::ChefSolo
    OS::Heat::CloudConfig
    OS::Heat::MultipartMime
    OS::Heat::None
    OS::Heat::RandomString
    OS::Heat::ResourceGroup
    OS::Heat::SoftwareConfig
    OS::Heat::SoftwareDeployment
    OS::Heat::SoftwareDeploymentGroup
    OS::Heat::SoftwareDeployments
    OS::Heat::Stack
    OS::Heat::SwiftSignal
    OS::Heat::SwiftSignalHandle
    OS::Neutron::Net
    OS::Neutron::Port
    OS::Neutron::SecurityGroup
    OS::Neutron::Subnet
    OS::Nova::KeyPair
    OS::Nova::Server
    OS::Swift::Container
    OS::Trove::Instance
    OS::Zaqar::Queue
    Rackspace::AutoScale::Group
    Rackspace::AutoScale::ScalingPolicy
    Rackspace::AutoScale::WebHook
    Rackspace::Cloud::BackupConfig
    Rackspace::Cloud::BigData
    Rackspace::Cloud::ChefSolo
    Rackspace::Cloud::CloudFilesCDN
    Rackspace::Cloud::DNS
    Rackspace::Cloud::LoadBalancer
    Rackspace::Cloud::Network
    Rackspace::Cloud::Server
    Rackspace::Cloud::WinServer
    Rackspace::CloudMonitoring::AgentToken
    Rackspace::CloudMonitoring::Alarm
    Rackspace::CloudMonitoring::Check
    Rackspace::CloudMonitoring::Entity
    Rackspace::CloudMonitoring::Notification
    Rackspace::CloudMonitoring::NotificationPlan
    Rackspace::CloudMonitoring::PlanNotifications
    Rackspace::Neutron::SecurityGroupAttachment
    Rackspace::RackConnect::PoolNode
    Rackspace::RackConnect::PublicIP

**Event**
~~~~~~~~~~~~

Event commands use this syntax::

    rack orchestration event <action> [command flags]

``get``
^^^^^^^^
Retrieves details for a specified event::

    rack orchestration event get --name <stackName> --resource <resourceName> --event <eventID> [optional flags]
    rack orchestration event get --id <stackID> --resource <resourceName> --event <eventID> [optional flags]

**Response**

.. code::

    $ rack orchestration event get --name RackTest --resource test_server --event dcfe8ad3-150f-4cbe-9993-2d82793753b7
    ResourceName		test_server
    Time			2015-09-13 04:20:24 +0000 UTC
    ResourceStatusReason	state changed
    ResourceStatus		ADOPT_COMPLETE
    PhysicalResourceID	f075a7c1-28ef-4699-9046-383098134902
    ID			dcfe8ad3-150f-4cbe-9993-2d82793753b7
    ResourceProperties	{
    			  "admin_pass": null,
    			  "admin_user": null,
    			  "availability_zone": null,
    			  "block_device_mapping": null,
    			  "block_device_mapping_v2": null,
    			  "config_drive": null,
    			  "diskConfig": null,
    			  "flavor": "2 GB General Purpose v1",
    			  "flavor_update_policy": "RESIZE",
    			  "image": "Debian 7 (Wheezy) (PVHVM)",
    			  "image_update_policy": "REBUILD",
    			  "key_name": null,
    			  "metadata": null,
    			  "name": "test-server-2",
    			  "networks": null,
    			  "personality": {},
    			  "reservation_id": null,
    			  "scheduler_hints": null,
    			  "security_groups": [],
    			  "software_config_transport": "POLL_TEMP_URL",
    			  "user_data": "",
    			  "user_data_format": "RAW"
    			}

``list-resource``
^^^^^^^^
Retrieves events for a specified stack resource::

    rack orchestration event list-resource --name <stackName> --resource <resourceName> [optional flags]
    rack orchestration event list-resource --id <stackID> --resource <resourceName> [optional flags]

**Response**

.. code::

    $ rack orchestration event list-resource --name RackTest --resource test_server
    ResourceName	Time				ResourceStatusReason	ResourceStatus		PhysicalResourceID			ID
    test_server	2015-09-13 04:20:24 +0000 UTC	state changed		ADOPT_COMPLETE		f075a7c1-28ef-4699-9046-383098134902	dcfe8ad3-150f-4cbe-9993-2d82793753b7
    test_server	2015-09-13 04:20:24 +0000 UTC	state changed		ADOPT_IN_PROGRESS						e78533e1-c8e0-4eca-8734-b193b6d32e06

``list-stack``
^^^^^^^^
Retrieves events for a specified stack resource::

    rack orchestration event list-stack --name <stackName> [optional flags]
    rack orchestration event list-resource --id <stackID> [optional flags]
    (echo stackName1 && echo stackName2) | rack orchestration event list-resource --stdin name [optional flags]

**Response**

.. code::

    $ rack orchestration event list-resource --name RackTest --resource test_server
    ResourceName	Time				ResourceStatusReason	ResourceStatus		PhysicalResourceID			ID
    test_server	2015-09-13 04:20:24 +0000 UTC	state changed		ADOPT_COMPLETE		f075a7c1-28ef-4699-9046-383098134902	dcfe8ad3-150f-4cbe-9993-2d82793753b7
    test_server	2015-09-13 04:20:24 +0000 UTC	state changed		ADOPT_IN_PROGRESS						e78533e1-c8e0-4eca-8734-b193b6d32e06

**Template**
~~~~~~~~~~~~

Template commands use this syntax::

    rack orchestration template <action> [command flags]

``validate``
^^^^^^^^
Validates a specified template::

    rack orchestration template validate --template <templateFile> [optional flags]
    rack orchestration template validate --template-url <templateURL> [optional flags]

**Response**

.. code::

    $ rack orchestration template validate --template-file my_template.yaml
    Description	No description
    Parameters	{
    		  "flavor": {
    		    "Default": 4353,
    		    "Description": "Flavor for the server to be created",
    		    "Label": "flavor",
    		    "NoEcho": "true",
    		    "Type": "String"
    		  }
    		}
    ParameterGroups	null

``get``
^^^^^^^
Retrieves template for a specified stack::

    rack orchestration template get --id <stackID> [optional flags]
    rack orchestration template get --name <stackName> [optional flags]
    (echo stackName1 && echo stackName2) | rack orchestration template get --stdin name

**Response**

.. code::

    $ rack orchestration template get --name RackTest
    {
      "heat_template_version": "2014-10-16",
      "parameters": {
        "flavor": {
          "default": 4353,
          "description": "Flavor for the server to be created",
          "hidden": true,
          "type": "string"
        }
      },
      "resources": {
        "test_server": {
          "properties": {
            "flavor": "2 GB General Purpose v1",
            "image": "Debian 7 (Wheezy) (PVHVM)",
            "name": "test-server-2"
          },
          "type": "OS::Nova::Server"
        }
      }
    }
