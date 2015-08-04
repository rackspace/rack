.. _files:

=====
Files
=====

This section lists all cloud files commands supported by Rackspace CLI.

Commands
--------

All cloud files commands are based on this syntax::

    rack files <resource> <action> [command flags]


``container``
~~~~~~~~~~~~~

All cloud files container operations use this syntax::

    rack files container <action> [optional flags]

``list``
^^^^^^^^
Retrieves a list of of containers::

    rack files container list [optional flags]

``create``
^^^^^^^^^^
Creates a container::

    rack files container create --name <containerName> [optional flags]
    (echo containerName1 && echo containerName2) | rack files container create --stdin name

``get``
^^^^^^^
Retrieves a list of containers. Optional flags can
be used to refine your search::

    rack files container get --name <containerName> [optional flags]
    (echo containerName1 && echo containerName2) | rack files container get --stdin name [optional flags]

``update``
^^^^^^^^^^
Creates or updates metadata for a specified container::

    rack files container update --name <containerName> [optional flags]

``delete``
^^^^^^^^^^
Permanently removes the specified container::

    rack files container delete --name <containerName> [optional flags]
    (echo containerName1 && echo containerName2) | rack files container delete --stdin name [optional flags]

``empty``
^^^^^^^^^
Empties a container of all its objects::

    rack files container empty --name <containerName> [optional flags]
    (echo containerName1 && echo containerName2) | rack files container empty --stdin name [optional flags]

``object``
~~~~~~~~~~

All cloud files object commands use this syntax::

    rack files object <action> [optional flags]

``list``
^^^^^^^^
Lists all objects contained in a specified container. Optional flags can be
used to refine your search::

    rack files object list --container <containerName> [optional flags]
    (echo containerName1 && echo containerName2) | rack files object list --stdin container [optional flags]

``upload``
^^^^^^^^^^
Uploads an object into a specified container::

    rack files object upload --container <containerName> --name <objectName> --content <contentToUpload> [optional flags]
    rack files object upload --container <containerName> --name <objectName> --file <fileToUpload> [optional flags]
    cat fileToUpload.txt | rack files object upload --container <containerName> --name <objectName> --stdin content [optional flags]

``upload-dir``
^^^^^^^^^^^^^^

Uploads an object directory into a specified container::

    rack files object upload-dir --container <containerName> --dir <dirPath> [optional flags]
    find . -type d -name "something*" | rack files object upload-dir --container <containerName> --stdin dir [optional flags]

``download``
^^^^^^^^^^^^
Downloads an object from the specified contained into your local system::

    rack files object download --container <containerName> --name <objectName> [optional flags]

``get``
^^^^^^^^
Retrieves an object's data::

    rack files object get --container <containerName> --name <objectName> [optional flags]

``delete``
^^^^^^^^^^
Permanently removes an object::

    rack files object delete --container <containerName> --name <objectName> [optional flags]
    (echo objectName1 && echo objectName2) | rack files object delete --container <containerName> --stdin name [optional flags]
