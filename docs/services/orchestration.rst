.. _orchestration:

======================
Orchestration commands
======================

This section lists the commands for interacting with Rackspace Cloud Orchestration.

All ``orchestration`` commands are based on the following syntax::

   rack orchestration <subcommand> <action> [command flags]

*Command flags* enable you to customize certain attributes of the command, such as using ``--name`` to name a stack. To display a list of command flags specific to the command, type ``rack orchestration <subcommand> <action> --help``.

The following sections describe the ``orchestration`` subcommands and the actions associated with them.

Build-info
----------

The ``build-info`` subcommand provides build information about the Cloud Orchestration service. The ``build-info`` subcommand uses the following syntax::

   rack orchestration build-info <action> [command flags]
   
The following section describes the action that you can perform on the ``build-info`` subcommand and provides an example response.

``get``
~~~~~~~

Retrieves the build information of the Orchestration service.

::

    rack orchestration build-info get [optional flags]

**Response**

.. code::

    $ rack orchestration build-info get
    API		2015.l3-20150903-1517
    Engine		2015.l3-20150903-1517
    FusionAPI	l1-20150622-17c7bae-141

Stack
-----

The ``stack`` subcommand provides information about and performs actions on the stacks in Cloud Orchestration. The ``stack`` subcommand uses the following syntax::

    rack orchestration stack <action> [command flags]

The following sections describe the actions that you can perform on the ``stack`` subcommand and provide example responses.

``list``
~~~~~~~~
Retrieves a list of stacks.

::

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
~~~~~~~~~~
Creates a stack with the specified name. You must select a template by using either the ``--template-file``  or ``--template-url`` flag to specify the file name or URL of the template that you want to use. 

::

    rack orchestration stack create --name <stackName> [optional flags]
    (echo stackName1 && echo stackName2) | rack orchestration stack create --stdin name [optional flags]

**Response**

.. code::

    $ rack orchestration stack create --name RackTest --template-file mytemplate.yaml
    rack orchestration stack create --name RackTest --template-file my_nova.yaml
    ID		27e5cb19-200d-4fbe-83a5-4f24fc4a3e9d
    Links0:Href	https://iad.orchestration.api.rackspacecloud.com/v1/TENANT_ID/stacks/RackTest/27e5cb19-200d-4fbe-83a5-4f24fc4a3e9d
    Links0:Rel	self

.. note::

    You can retrieve the details of the stack that you created by using the ``get`` command described in the following section.

``get``
~~~~~~~
Retrieves the details about a stack, which you can specify by ID or name.

::

    rack orchestration stack get --id <stackID> [optional flags]
    rack orchestration stack get --name <stackName> [optional flags]
    (echo stackName1 && echo stackName2) | rack orchestration stack get --stdin name [optional flags]

**Response**

.. code::

    $ rack orchestration stack get --name RackTest
    CreationTime			2015-10-15 19:51:43 +0000 UTC
    Description			No description
    DisableRollback			true
    ID				27e5cb19-200d-4fbe-83a5-4f24fc4a3e9d
    Name				RackTest
    Status				CREATE_COMPLETE
    StatusReason			Stack CREATE completed successfully
    Tags				[]
    TemplateDescription		No description
    Timeout				0
    UpdatedTime
    Parameters:OS::stack_name	RackTest
    Parameters:OS::project_id	TENANT_ID
    Parameters:flavor		******
    Parameters:OS::stack_id		27e5cb19-200d-4fbe-83a5-4f24fc4a3e9d
    Links0:Href			https://iad.orchestration.api.rackspacecloud.com/v1/TENANT_ID/stacks/RackTest/27e5cb19-200d-4fbe-83a5-4f24fc4a3e9d
    Links0:Rel			self

``update``
~~~~~~~~~~
Updates a stack by using a provided template. You can specify the stack by ID or name, and you can specify the template by file name or URL.

::

    rack orchestration stack update --id <stackID> [optional flags]
    rack orchestration stack update --name <stackName> [optional flags]

**Response**

.. code::

    $ rack orchestration stack update --name RackTest --template-file myNewTemplate.yaml
    CreationTime			2015-10-15 19:51:43 +0000 UTC
    Description			No description
    DisableRollback			true
    ID				27e5cb19-200d-4fbe-83a5-4f24fc4a3e9d
    Name				RackTest
    Status				UPDATE_IN_PROGRESS
    StatusReason			Stack UPDATE started
    Tags				[]
    TemplateDescription		No description
    Timeout				0
    UpdatedTime
    Parameters:flavor		******
    Parameters:OS::stack_id		27e5cb19-200d-4fbe-83a5-4f24fc4a3e9d
    Parameters:OS::stack_name	RackTest
    Parameters:OS::project_id	TENANT_ID
    Links0:Href			https://iad.orchestration.api.rackspacecloud.com/v1/TENANT_ID/stacks/RackTest/27e5cb19-200d-4fbe-83a5-4f24fc4a3e9d
    Links0:Rel			self

``delete``
~~~~~~~~~~
Deletes a stack, which you can specify by ID or name. 

::

    rack orchestration stack delete --id <stackID> [optional flags]
    rack orchestration stack delete --name <stackName> [optional flags]
    (echo stackName1 && echo stackName2) | rack orchestration stack delete --stdin name [optional flags]

**Response**

.. code::

    $ rack orchestration stack delete --name RackTest
    Stack RackTest is being deleted.

``preview``
~~~~~~~~~~~
Shows the number and type of resources that will be created in a stack by the specified template.

::

    rack  orchestration stack preview --name <stackName> [--template-file <templateFile> | --template-url <templateURL>] [optional flags]
    (echo stackName1 && echo stackName2) | rack  orchestration stack preview --stdin name [--template-file <templateFile> | --template-url <templateURL>] [optional flags]

**Response**

.. code::

    $ rack orchestration stack preview --template-file my_template.yaml --name RackTest
    CreationTime					2015-10-15 20:42:42.862279 +0000 UTC
    Description					No description
    DisableRollback					true
    ID						None
    Name						RackTest2
    TemplateDescription				No description
    Timeout						0
    UpdatedTime
    Links0:Href					https://iad.orchestration.api.rackspacecloud.com/v1/TENANT_ID/stacks/RackTest2/None
    Links0:Rel					self
    Parameters:OS::project_id			TENANT_ID
    Parameters:flavor				******
    Parameters:OS::stack_id				None
    Parameters:OS::stack_name			RackTest2
    Resources0:resource_identity:stack_name		RackTest2
    Resources0:resource_identity:stack_id		None
    Resources0:resource_identity:tenant		TENANT_ID
    Resources0:resource_identity:path		/resources/test_server
    Resources0:resource_action			INIT
    Resources0:attributes:accessIPv4
    Resources0:attributes:accessIPv6
    Resources0:attributes:networks
    Resources0:attributes:addresses
    Resources0:attributes:console_urls
    Resources0:attributes:name			test-server
    Resources0:attributes:first_address
    Resources0:attributes:instance_name
    Resources0:resource_name			test_server
    Resources0:stack_identity:stack_name		RackTest2
    Resources0:stack_identity:stack_id		None
    Resources0:stack_identity:tenant		TENANT_ID
    Resources0:stack_identity:path
    Resources0:stack_name				RackTest2
    Resources0:resource_status			COMPLETE
    Resources0:updated_time
    Resources0:properties:key_name
    Resources0:properties:config_drive
    Resources0:properties:availability_zone
    Resources0:properties:image			Debian 7 (Wheezy) (PVHVM)
    Resources0:properties:reservation_id
    Resources0:properties:flavor_update_policy	RESIZE
    Resources0:properties:flavor			2 GB General Purpose v1
    Resources0:properties:networks0:port
    Resources0:properties:networks0:subnet
    Resources0:properties:networks0:fixed_ip
    Resources0:properties:networks0:uuid
    Resources0:properties:networks0:network		11111111-1111-1111-1111-111111111111
    Resources0:properties:user_data_format		RAW
    Resources0:properties:admin_user
    Resources0:properties:name			test-server
    Resources0:properties:user_data
    Resources0:properties:diskConfig
    Resources0:properties:scheduler_hints
    Resources0:properties:metadata
    Resources0:properties:block_device_mapping_v2
    Resources0:properties:admin_pass
    Resources0:properties:block_device_mapping
    Resources0:properties:software_config_transport	POLL_TEMP_URL
    Resources0:properties:image_update_policy	REBUILD
    Resources0:description
    Resources0:creation_time
    Resources0:resource_type			OS::Nova::Server
    Resources0:resource_status_reason
    Resources0:physical_resource_id

``abandon``
~~~~~~~~~~~
Abandons a stack, which deletes the record of the stack from Orchestration but does not delete any of the underlying resources. You can specify the stack by ID or name.

::

    rack orchestration stack abandon --id <stackID> [optional flags]
    rack orchestration stack abandon --name <stackName> [optional flags]

To obtain a JSON representation of the abandoned stack, use the ``--output json`` flag. When this JSON is stored in a file, you can use it in the ``adopt`` command to create a new stack with the resources of the abandoned stack.

**Response**

.. code::

    $ rack orchestration stack abandon --name RackTest
    Status							COMPLETE
    Name							RackTest
    Action							CREATE
    ID							22e669f3-510f-4ef1-8782-96ad692d8b41
    StackUserProjectID					TENANT_ID
    ProjectID						TENANT_ID
    Template:heat_template_version				2014-10-16
    Template:resources:test_server:properties:flavor	2 GB General Purpose v1
    Template:resources:test_server:properties:name		test-server
    Template:resources:test_server:properties:image		Debian 7 (Wheezy) (PVHVM)
    Template:resources:test_server:properties:networks0:uuid11111111-1111-1111-1111-111111111111
    Template:resources:test_server:type			OS::Nova::Server
    Template:parameters:flavor:type				string
    Template:parameters:flavor:description			Flavor for the server to be created
    Template:parameters:flavor:default			4353
    Template:parameters:flavor:hidden			true
    Resources:test_server:status				COMPLETE
    Resources:test_server:name				test_server
    Resources:test_server:resource_id			e3a5c760-25fc-4a96-915d-a3dcbf94019a
    Resources:test_server:action				CREATE
    Resources:test_server:type				OS::Nova::Server

``adopt``
~~~~~~~~~
Creates a stack without creating any resources; existing resources are used instead. This command is usually used to create a stack by using the resources of an abandoned stack. You can use the JSON output representation of the abandoned stack as the contents for ``adopt-file`` to direct Orchestration to use the resources of the abandoned stack in the creation of the adopted stack.

::

    rack orchestration stack adopt --name stackName --adopt-file adoptFile [optional flags]

**Response**

.. code::

    $ rack orchestration stack adopt --name RackTest --adopt-file abandon.yaml
    ID		27e5cb19-200d-4fbe-83a5-4f24fc4a3e9d
    Links0:Href	https://iad.orchestration.api.rackspacecloud.com/v1/TENANT_ID/stacks/RackTest/27e5cb19-200d-4fbe-83a5-4f24fc4a3e9d
    Links0:Rel	self

``list-events``
~~~~~~~~~~~~~~~
Retrieves events for a specified stack, which you can specify by ID or name. 

::

    rack orchestration stack list-events --name <stackName> [optional flags]
    rack orchestration stack list-events --id <stackID> [optional flags]
    (echo stackName1 && echo stackName2) | rack orchestration stack list-events --stdin name [optional flags]

**Response**

.. code::

    $ rack orchestration stack list-events --name RackTest --resource-names test_server
    ResourceName	Time				ResourceStatusReason	ResourceStatus		PhysicalResourceID			ID
    test_server	2015-09-13 04:20:24 +0000 UTC	state changed		ADOPT_COMPLETE		f075a7c1-28ef-4699-9046-383098134902	dcfe8ad3-150f-4cbe-9993-2d82793753b7
    test_server	2015-09-13 04:20:24 +0000 UTC	state changed		ADOPT_IN_PROGRESS						e78533e1-c8e0-4eca-8734-b193b6d32e06

``get-template``
~~~~~~~~~~~~~~~~
Retrieves the template for a stack, which you can specify by ID or name.

::

    rack orchestration stack get-template --id <stackID> [optional flags]
    rack orchestration stack get-template --name <stackName> [optional flags]
    (echo stackName1 && echo stackName2) | rack orchestration stack get-template --stdin name

**Response**

.. code::

    $ rack orchestration stack get-template --name RackTest
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


Resource
--------

The ``resource`` subcommand provides information about and performs actions on the resources in Cloud Orchestration. The ``resource`` subcommand uses the following syntax::

    rack orchestration resource <action> [command flags]

The following sections describe the actions that you can perform on the ``resource`` subcommand and provide example responses.

``list``
~~~~~~~~
Retrieves a list of resources for a given stack, which you can specify by ID or name. 

::

    rack orchestration resource list --stack-name <stackName> [optional flags]
    rack orchestration resource list --stack-id <stackID> [optional flags]
    (echo stackName1 && echo stackName2)  | rack orchestration resource list --stdin stack-name [optional flags]

**Response**

.. code::

    $ rack orchestration resource list --stack-name RackTest
    Name		PhysicalID				                Type			    Status		               UpdatedTime
    test_server	f075a7c1-28ef-4699-9046-383098134902	OS::Nova::Server	CREATE_COMPLETE	2015-09-12 16:37:49 +0000 UTC

``get``
~~~~~~~
Retrieves the details about a specified resource in a stack, which you can specify by ID or name. 

::

    rack orchestration resource get --stack-id <stackID> --name <resourceName> [optional flags]
    rack orchestration resource get --stack-name <stackName> --name <resourceName> [optional flags]

**Response**

.. code::

    $ rack orchestration resource get --stack-name RackTest --name test_server
    CreationTime					2015-10-15 21:09:54 +0000 UTC
    Description
    LogicalID					test_server
    Name						test_server
    PhysicalID					d4ffd1fd-ef29-4e31-8776-78414a5c0f67
    Status						CREATE_COMPLETE
    StatusReason					state changed
    Type						OS::Nova::Server
    UpdatedTime					2015-10-15 21:09:54 +0000 UTC
    Attributes:OS-EXT-STS:power_state		1
    Attributes:metadata:rax_service_level_automationIn Progress
    Attributes:image:links0:href			https://iad.servers.api.rackspacecloud.com/TENANT_ID/images/eaaf94d8-55a6-4bfa-b0a8-473febb012dc
    Attributes:image:links0:rel			bookmark
    Attributes:image:id				eaaf94d8-55a6-4bfa-b0a8-473febb012dc
    Attributes:OS-DCF:diskConfig			MANUAL
    Attributes:id					d4ffd1fd-ef29-4e31-8776-78414a5c0f67
    Attributes:OS-EXT-STS:vm_state			active
    Attributes:user_id				5c11b69d82cf4313b7a8b173b799a0ef
    Attributes:tenant_id				TENANT_ID
    Attributes:accessIPv4
    Attributes:created				2015-10-15T21:09:56Z
    Attributes:addresses:private0:version		4
    Attributes:addresses:private0:addr		10.208.234.120
    Attributes:accessIPv6
    Attributes:status				ACTIVE
    Attributes:RAX-PUBLIC-IP-ZONE-ID:publicIPZoneId	025e96cf138a9036fffb45031c506ac7a7052a355b7f08bcbbc12da9
    Attributes:flavor:id				general1-2
    Attributes:flavor:links0:href			https://iad.servers.api.rackspacecloud.com/TENANT_ID/flavors/general1-2
    Attributes:flavor:links0:rel			bookmark
    Attributes:links0:rel				self
    Attributes:links0:href				https://iad.servers.api.rackspacecloud.com/v2/TENANT_ID/servers/d4ffd1fd-ef29-4e31-8776-78414a5c0f67
    Attributes:links1:href				https://iad.servers.api.rackspacecloud.com/TENANT_ID/servers/d4ffd1fd-ef29-4e31-8776-78414a5c0f67
    Attributes:links1:rel				bookmark
    Attributes:key_name
    Attributes:OS-EXT-STS:task_state
    Attributes:progress				100
    Attributes:name					test-server
    Attributes:hostId				c1529238ae34923ed243a257ffb72e92db13ab2552994f76b26f3ce7
    Attributes:config_drive
    Attributes:updated				2015-10-15T21:10:36Z
    Links0:Href					https://iad.orchestration.api.rackspacecloud.com/v1/TENANT_ID/stacks/RackTest/deb6e034-2808-4db6-9807-fa00e9709925/resources/test_server
    Links0:Rel					self
    Links1:Href					https://iad.orchestration.api.rackspacecloud.com/v1/TENANT_ID/stacks/RackTest/deb6e034-2808-4db6-9807-fa00e9709925
    Links1:Rel					stack

``get-schema``
~~~~~~~~~~~~~~
Shows the interface schema for a specified resource type.

::

    rack orchestration resource get-schema --type <resourceType> [optional flags]
    (echo resourceType1 && echo resourceType2) | rack orchestration resource get-schema --stdin type [optional flags]

This schema describes the properties that can be set on the resource, their types, constraints, descriptions, and default values. Additionally, the resource attributes and their descriptions are provided.

**Response**

.. code::

    $ rack orchestration resource get-schema --type OS::Heat::None
    ResourceType			OS::Heat::None
    Attributes:show:type		map
    Attributes:show:description	Detailed information about resource.
    SupportStatus:status		SUPPORTED
    SupportStatus:message
    SupportStatus:version		5.0.0
    SupportStatus:previous_status

``get-template``
~~~~~~~~~~~~~~~~
Shows a template representation for the specified resource type.

::

    rack orchestration resource get-template --type <resourceType> [optional flags]
    (echo resourceType1 && echo resourceType2) | rack orchestration resource get-template --stdin type [optional flags]

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
~~~~~~~~~~~~~~
Lists the supported template resource types.

::

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

``list-events``
~~~~~~~~~~~~~~~
Retrieves events for a specified stack resource.

::

    rack orchestration resource list-events --stack-name <stackName> --name <resourceName> [optional flags]
    rack orchestration resource list-events --stack-id <stackID> --name <resourceName> [optional flags]

**Response**

.. code::

    $ rack orchestration resource list-events --stack-name RackTest --resource-name test_server
    ResourceName	Time				ResourceStatusReason	ResourceStatus		PhysicalResourceID			ID
    test_server	2015-09-13 04:20:24 +0000 UTC	state changed		ADOPT_COMPLETE		f075a7c1-28ef-4699-9046-383098134902	dcfe8ad3-150f-4cbe-9993-2d82793753b7
    test_server	2015-09-13 04:20:24 +0000 UTC	state changed		ADOPT_IN_PROGRESS						e78533e1-c8e0-4eca-8734-b193b6d32e06

Event
-----

The ``event`` subcommand provides information about events in Cloud Orchestration. The ``event`` subcommand uses the following syntax::

    rack orchestration event <action> [command flags]

The following section describes the action that you can perform on the ``event`` subcommand and provides an example response.

``get``
~~~~~~~
Retrieves details for a specified event.

::

    rack orchestration event get --stack-name <stackName> --resource <resourceName> --id <eventID> [optional flags]
    rack orchestration event get --stack-id <stackID> --resource <resourceName> --id <eventID> [optional flags]

**Response**

.. code::

    $ rack orchestration event get --stack-name RackTest --resource test_server --id c4beb7be-8c8a-4e6a-ad52-b5c571564c77
    ResourceName					test_server
    Time						2015-10-15 21:10:20 +0000 UTC
    ResourceStatusReason				state changed
    LogicalResourceID				test_server
    ResourceStatusReason				state changed
    ResourceStatus					CREATE_COMPLETE
    PhysicalResourceID				d4ffd1fd-ef29-4e31-8776-78414a5c0f67
    ID						c4beb7be-8c8a-4e6a-ad52-b5c571564c77
    Links0:Href					https://iad.orchestration.api.rackspacecloud.com/v1/TENANT_ID/stacks/RackTest/deb6e034-2808-4db6-9807-fa00e9709925/resources/test_server/events/c4beb7be-8c8a-4e6a-ad52-b5c571564c77
    Links0:Rel					self
    Links1:Href					https://iad.orchestration.api.rackspacecloud.com/v1/TENANT_ID/stacks/RackTest/deb6e034-2808-4db6-9807-fa00e9709925/resources/test_server
    Links1:Rel					resource
    Links2:Href					https://iad.orchestration.api.rackspacecloud.com/v1/TENANT_ID/stacks/RackTest/deb6e034-2808-4db6-9807-fa00e9709925
    Links2:Rel					stack
    ResourceProperties:networks0:fixed_ip
    ResourceProperties:networks0:port
    ResourceProperties:networks0:network		11111111-1111-1111-1111-111111111111
    ResourceProperties:networks0:uuid
    ResourceProperties:networks0:subnet
    ResourceProperties:metadata
    ResourceProperties:user_data_format		RAW
    ResourceProperties:admin_pass
    ResourceProperties:flavor_update_policy		RESIZE
    ResourceProperties:diskConfig
    ResourceProperties:flavor			2 GB General Purpose v1
    ResourceProperties:config_drive
    ResourceProperties:reservation_id
    ResourceProperties:key_name
    ResourceProperties:scheduler_hints
    ResourceProperties:block_device_mapping
    ResourceProperties:block_device_mapping_v2
    ResourceProperties:name				test-server
    ResourceProperties:software_config_transport	POLL_TEMP_URL
    ResourceProperties:user_data
    ResourceProperties:admin_user
    ResourceProperties:image_update_policy		REBUILD
    ResourceProperties:availability_zone
    ResourceProperties:image			Debian 7 (Wheezy) (PVHVM)

Template
--------

The ``template`` subcommand provides information about and performs actions on the templates in Cloud Orchestration. The ``template`` subcommand uses the following syntax::

    rack orchestration template <action> [command flags]

The following section describes the action that you can perform on the ``template`` subcommand and provides an example response.

``validate``
~~~~~~~~~~~~
Validates a specified template.

::

    rack orchestration template validate --template <templateFile> [optional flags]
    rack orchestration template validate --template-url <templateURL> [optional flags]

**Response**

.. code::

    $ rack orchestration template validate --template-file my_template.yaml
    Description			No description
    Parameters:flavor:Default	4353
    Parameters:flavor:NoEcho	true
    Parameters:flavor:Type		String
    Parameters:flavor:Description	Flavor for the server to be created
    Parameters:flavor:Label		flavor
