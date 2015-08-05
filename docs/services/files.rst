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

Cloud files container commands use this syntax::

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
Retrieves a list of containers::

    rack files container get --name <containerName> [optional flags]
    (echo containerName1 && echo containerName2) | rack files container get --stdin name [optional flags]

``update``
^^^^^^^^^^
Create or update read and write permissions for a specified container::

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


``set-metadata``
^^^^^^^^^^^^^^^^
Sets metadata for the specified container::

    rack files container set-metadata --name <conatinerName> --metadata <key1=val1,key2=val2,...> [optional flags]

``get-metadata``
^^^^^^^^^^^^^^^^
Retrieves the metadata for a given container::

    rack files container get-metadata --name <conatinerName> [optional flags]

``update-metadata``
^^^^^^^^^^^^^^^^^^^
Updates metadata items for a specified container, or adds the specified
metadata if there is no current metadata associated with the container::

    rack files container update-metadata --name <conatinerName> --metadata <key1=val1,key2=val2,...> [optional flags]

``delete-metadata``
^^^^^^^^^^^^^^^^^^^
Deletes one or more metadata keys for a container::

    rack files container delete-metadata --name <conatinerName> --metadata-keys <key1,key2,...> [optional flags]


``object``
~~~~~~~~~~

Cloud files object commands use this syntax::

    rack files object <action> [optional flags]

``list``
^^^^^^^^
Lists all objects contained in a specified container::

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

``set-metadata``
^^^^^^^^^^^^^^^^
Sets metadata for the specified object::

    rack files object set-metadata --container <containerName> --name <objectName> --metadata <key1=val1,key2=val2,...> [optional flags]

``get-metadata``
^^^^^^^^^^^^^^^^
Retrieves the metadata for a given object::

    rack files object get-metadata --container <containerName> --name <objectName> [optional flags]

``update-metadata``
^^^^^^^^^^^^^^^^^^^
Updates metadata items for a specified object, or adds the specified
metadata if there is no current metadata associated with the object::

    rack files object update-metadata --container <containerName> --name <objectName> --metadata <key1=val1,key2=val2,...> [optional flags]

``delete-metadata``
^^^^^^^^^^^^^^^^^^^
Deletes one or more metadata keys from an object::

    rack files object delete-metadata --container <containerName> --name <objectName> --metadata-keys <key1,key2,...> [optional flags]


``account``
~~~~~~~~~~

Cloud Files account commands use this syntax::

    rack files account <action> [optional flags]

``set-metadata``
^^^^^^^^^^^^^^^^
Sets metadata for the specified account::

    rack files account set-metadata --metadata <key1=val1,key2=val2,...> [optional flags]

``get-metadata``
^^^^^^^^^^^^^^^^
Retrieves the metadata for a given account::

    rack files account get-metadata [optional flags]

``update-metadata``
^^^^^^^^^^^^^^^^^^^
Updates metadata items for a specified account, or adds the specified
metadata if there is no current metadata associated with the account::

    rack files account update-metadata --metadata <key1=val1,key2=val2,...> [optional flags]

``delete-metadata``
^^^^^^^^^^^^^^^^^^^
Deletes one or more metadata keys from an account::

    rack files account delete-metadata  --metadata-keys <key1,key2,...> [optional flags]

