.. _files:

files
=======

Commands for Rackspace Cloud Files.

Synopsis
--------

::

   rack files <resource> <action> [command flags]

Commands
--------

``container``
~~~~~~~~~~~~~

  Files Container operations

``list``
^^^^^^^^
Usage::

    rack files container list [optional flags]

``create``
^^^^^^^^^^
Usage::

    rack files container create --name <containerName> [optional flags]
    (echo containerName1 && echo containerName2) | rack files container create --stdin name

``get``
^^^^^^^
Usage::

    rack files container get --name <containerName> [optional flags]
    (echo containerName1 && echo containerName2) | rack files container get --stdin name [optional flags]

``update``
^^^^^^^^^^
Usage::

    rack files container update --name <containerName> [optional flags]

``delete``
^^^^^^^^^^
Usage::

    rack files container delete --name <containerName> [optional flags]
    (echo containerName1 && echo containerName2) | rack files container delete --stdin name [optional flags]

``empty``
^^^^^^^^^
Usage::

    rack files container empty --name <containerName> [optional flags]
    (echo containerName1 && echo containerName2) | rack files container empty --stdin name [optional flags]

``set-metadata``
^^^^^^^^^^^^^^^^
Usage::

    rack files container set-metadata --name <containerName> --metadata <metadata> [optional flags]

``get-metadata``
^^^^^^^^^^^^^^^^
Usage::

    rack files container get-metadata --name <containerName> [optional flags]

``update-metadata``
^^^^^^^^^^^^^^^^
Usage::

    rack files container update-metadata --name <containerName> --metadata <metadata> [optional flags]

``delete-metadata``
^^^^^^^^^^^^^^^^
Usage::

    rack files container delete-metadata --name <containerName> --metadata-keys <metadataKeys> [optional flags]


``object``
~~~~~~~~~~

  Files Object operations

``list``
^^^^^^^^
Usage::

    rack files object list --container <containerName> [optional flags]
    (echo containerName1 && echo containerName2) | rack files object list --stdin container [optional flags]

``upload``
^^^^^^^^^^
Usage::

    rack files object upload --container <containerName> --name <objectName> --content <contentToUpload> [optional flags]
    rack files object upload --container <containerName> --name <objectName> --file <fileToUpload> [optional flags]
    cat fileToUpload.txt | rack files object upload --container <containerName> --name <objectName> --stdin content [optional flags]

``upload-dir``
^^^^^^^^^^^^^^

Usage::

    rack files object upload-dir --container <containerName> --dir <dirPath> [optional flags]
    find . -type d -name "something*" | rack files object upload-dir --container <containerName> --stdin dir [optional flags]

``download``
^^^^^^^^^^^^
Usage::

    rack files object download --container <containerName> --name <objectName> [optional flags]

``get``
^^^^^^^^
Usage::

    rack files object get --container <containerName> --name <objectName> [optional flags]

``delete``
^^^^^^^^^^
Usage::

    rack files object delete --container <containerName> --name <objectName> [optional flags]
    (echo objectName1 && echo objectName2) | rack files object delete --container <containerName> --stdin name [optional flags]

``set-metadata``
^^^^^^^^^^^^^^^^
Usage::

    rack files object set-metadata --name <objectName> --container <containerName> --metadata <metadata> [optional flags]

``get-metadata``
^^^^^^^^^^^^^^^^
Usage::

    rack files object get-metadata --name <objectName> --container <containerName> [optional flags]

``update-metadata``
^^^^^^^^^^^^^^^^
Usage::

    rack files object update-metadata --name <objectName> --container <containerName> --metadata <metadata> [optional flags]

``delete-metadata``
^^^^^^^^^^^^^^^^
Usage::

    rack files object delete-metadata --name <objectName> --container <containerName> --metadata-keys <metadataKeys> [optional flags]
