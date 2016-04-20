.. _cloudfilesexamples:

==============
Files examples
==============

Before you start using examples, be sure to enter your username and API key and store them locally by running the ``rack configure`` command. For more information, see :ref:`installation_and_configuration`.

You can get help for any command and its options by adding ``--help`` to the series of commands::

    $ rack files container create --help

Search for existing objects
~~~~~~~~~~~~~~~~~~~~~~~~~~~

If you have a large number of objects, you can search for the ones that you want.

First, list the available containers in your account.

::

    $ rack files container list
    Name						Count	Bytes
    Presentations				 14	    125902668
    Public						 47	    179289516
    allphotos					 493	295697436

Then specify a container and search for a specific term. The following example searches through a container named ``Presentations`` for a file with the name ``workshop`` in it. Use the appropriate search command for your operating system. 

::

    $ rack files object list --container Presentations | grep "workshop"
    humanitarian-openstack-workshop.zip     16760097    application/zip    2014-09-17T03:18:11.873700

Delete all objects in a container with a matching pattern
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Use this command as a convenient one-line cleanup utility: list the objects in a container, use a pattern to find certain objects, and then delete all of those objects. The following example shows clearing all the ``.log`` files from a container. Use the appropriate search command for your operating system.

::

    $ rack files object list --all-pages --container server-test-logs --fields name | \
    grep -i '.log' | \
    rack files object delete --container server-test-logs --stdin name

Download an object
~~~~~~~~~~~~~~~~~~

If you know the name of an object and the container it is stored in, you can download the object.

::

    $ rack files object download --container Presentations --name \
    humanitarian-openstack-workshop.zip > humanitarian-openstack-workshop.zip

Upload an object
~~~~~~~~~~~~~~~~

When you want to upload an object to Cloud Files, first ensure that you have a container to upload to.

::

    $ rack files container create --name screenshots

Next, upload the object, which in the following example is a screenshot.

::

    $ rack files object upload --container screenshots --file browser-screenshot.png --name browser-screenshot.png
