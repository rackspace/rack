.. _orchestrationexamples:

======================
Orchestration examples
======================

Before you start using examples, be sure to enter your username and API key and store them locally by running the ``rack configure`` command. For more information, see :ref:`installation_and_configuration`.

You can get help for any command and its options by appending --help to the series of commands::

    $ rack orchestration stack create --help

Create a stack
~~~~~~~~~~~~~~

Start by creating a template for the stack. Create a new file named ``teststack.yaml`` and add the following content to it.

::

    heat_template_version: 2013-05-23
    
    description: |
      Template for the Heat Smoke Test
    
    resources:
    
      random_key_name:
        type: OS::Heat::RandomString
        properties:
          length: 8

Then, use the following command to create a stack from the template.

::

    $ rack orchestration stack create --name RackTest --template-file teststack.yaml

    ID		cb10ba04-cf2d-479d-8355-ede637b3d487
    Links0:Href	https://dfw.orchestration.api.rackspacecloud.com/v1/123456/stacks/RackTest/cb10ba04-cf2d-479d-8355-ede637b3d487
    Links0:Rel	self

Abandon a stack
~~~~~~~~~~~~~~~

If you want to delete a stack without deleting the resources in it, you can *abandon* the stack. If you specify JSON output and save the output to a file, you can create (adopt) a stack later that uses those same resources.

::

    $ rack orchestration stack abandon --name RackTest --output json >> teststackabandon.json

Adopt a stack
~~~~~~~~~~~~~

If you want to create a stack by using the resources from an abandoned stack, you can *adopt* it from the stored output. 

::

    $ rack orchestration stack adopt --name RackTest --adopt-file teststackabandon.json

    ID		9da703de-d520-4088-ac8f-03c6d7cbd9ba
    Links0:Href	https://dfw.orchestration.api.rackspacecloud.com/v1/123456/stacks/RackTest/9da703de-d520-4088-ac8f-03c6d7cbd9ba
    Links0:Rel	self
