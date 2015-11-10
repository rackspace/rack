.. _orchestrationexamples:

=============
Orchestration
=============

Before you get started on any examples, be sure you have entered your
username and API key and stored them locally::

    rack configure

Create a stack
~~~~~~~~~~~~~~

Let's start out by creating a stack. Create a new file named ``teststack.yaml``
and add the following to it::

    heat_template_version: 2013-05-23
    
    description: |
      Template for the Heat Smoke Test
    
    resources:
    
      random_key_name:
        type: OS::Heat::RandomString
        properties:
          length: 8

Then, use ``rack`` to create a stack from it::

    $ rack orchestration stack create --name RackTest --template-file teststack.yaml

    ID		cb10ba04-cf2d-479d-8355-ede637b3d487
    Links0:Href	https://dfw.orchestration.api.rackspacecloud.com/v1/123456/stacks/RackTest/cb10ba04-cf2d-479d-8355-ede637b3d487
    Links0:Rel	self

Abandon a stack
~~~~~~~~~~~~~~~

If we want to delete the stack without deleting the resources, we can ``abandon`` it while asking for
JSON output and saving the output to a file::

    $ rack orchestration stack abandon --name RackTest --output json >> teststackabandon.json

Adopt a stack
~~~~~~~~~~~~~

Now if we later want to create the stack again with those same resources, we can ``adopt`` it from the stored output::

    $ rack orchestration stack adopt --name RackTest --adopt-file teststackabandon.json

    ID		9da703de-d520-4088-ac8f-03c6d7cbd9ba
    Links0:Href	https://dfw.orchestration.api.rackspacecloud.com/v1/123456/stacks/RackTest/9da703de-d520-4088-ac8f-03c6d7cbd9ba
    Links0:Rel	self
