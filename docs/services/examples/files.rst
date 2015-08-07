.. _cloudfilesexamples:

====================
Files
====================

Before you get started on any examples, be sure you have entered your
username and API key and stored them locally::

    rack configure

Search for existing objects
~~~~~~~~~~~~~~~~~~~~~~~~~~~

If you have a lot of objects, the `rack` command lets you search through
the list.

First, list the available containers in your account::

    $ rack files container list

    Name						Count	Bytes
    Presentations					14	125902668
    Public						47	179289516
    allphotos					493	295697436

Here's an example of searching through a container named
Presentations for a file with the name "workshop" in it::

    $ rack files object list --container Presentations | grep "workshop"

    humanitarian-openstack-workshop.zip     16760097    application/zip    2014-09-17T03:18:11.873700

Delete all objects in a container with a matching pattern
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

As a handy one-line clean-up utility, use this command to list objects in a
container, then match on a pattern, then delete them all. This example shows
clearing all the .log files from a container::

    rack files object list --all-pages --container server-test-logs --fields name | \
    grep -i '.log' | \
    rack files object delete --container server-test-logs --stdin name

Download an object
~~~~~~~~~~~~~~~~~~

With the name of an object and the container it is stored within you can download it::

    rack files object download --container Presentations --name \
    humanitarian-openstack-workshop.zip > humanitarian-openstack-workshop.zip

Upload an object
~~~~~~~~~~~~~~~~

First, make sure you have a container to upload to::

    rack files container create --name screenshots

Next, upload a file, in this case a screenshot::

    rack files object upload --container screenshots --file browser-screenshot.png --name browser-screenshot.png
