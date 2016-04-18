.. _files:

==============
Files commands
==============

This section lists the commands for interacting with Rackspace Cloud Files.

All ``files`` commands are based on the following syntax::

    rack files <subcommand> <action> [command flags]

*Command flags* enable you to customize certain attributes of the command, such as using ``--name`` to name a Cloud Files container. To display a list of command flags specific to the command, type ``rack files <subcommand> <action> --help``.

The following sections describe the ``files`` subcommands and the actions associated with them.

Container
---------

The ``container`` subcommand provides information about and performs actions on the containers in Cloud Files. The ``container`` subcommand uses the following syntax::

    rack files container <action> [optional flags]

The following sections describe the actions that you can perform on the ``container`` subcommand and provide example responses.

``list``
~~~~~~~~
Retrieves a list of containers.

::

    rack files container list [optional flags]

**Response**

.. code::

    $ rack files container list
    Name					                        Count	Bytes
    770fb26f-2c43-4196-95d1-ad9ec1008365	2	    32
    Screenshots				                    5	    343148
    Work Docs			                       	4	    674447
    cloudservers			                  	45	  2296934481
    my test container			                0	    0

``create``
~~~~~~~~~~
Creates a container.

::

    rack files container create --name <containerName> [optional flags]
    (echo containerName1 && echo containerName2) | rack files container create --stdin name

**Response**

.. code::

    $ rack files container create --name RackCLI
    Successfully created container [RackCLI]

``get``
~~~~~~~
Retrieves the details for a specified container.

::

    rack files container get --name <containerName> [optional flags]
    (echo containerName1 && echo containerName2) | rack files container get --stdin name [optional flags]

**Response**

.. code::

    $ rack files container get --name RackCLI
    Name		RackCLI
    ObjectCount	0
    BytesUsed	0
    ContentLength	0
    AcceptRanges	bytes
    ContentType	application/json; charset=utf-8
    Read
    Write
    TransID		txc108b4d4ad9040898a210-0055c8d7dfdfw1
    VersionsLocation

``update``
~~~~~~~~~~
Creates or updates read and write permissions for a specified container.

::

    rack files container update --name <containerName> [optional flags]

**Response**

.. code::

    $ rack files container update --name RackCLI --container-read user1
    Successfully updated container [RackCLI]

``delete``
~~~~~~~~~~
Permanently deletes the specified container.

::

    rack files container delete --name <containerName> [optional flags]
    (echo containerName1 && echo containerName2) | rack files container delete --stdin name [optional flags
    
**Response**

.. code::

    $ rack files container delete --name RackCLI
    Running containers.Delete
    Successfully deleted container [RackCLI]

``empty``
~~~~~~~~~
Empties a container of all its objects.

::

    rack files container empty --name <containerName> [optional flags]
    (echo containerName1 && echo containerName2) | rack files container empty --stdin name [optional flags]

**Response**

.. code::

    $ rack files container empty --name RackCLI
    Successfully deleted object [colors-logos.png] from container [RackCLI]

    Successfully deleted object [dashed-lines.png] from container [RackCLI]

    Successfully deleted object [arrowsscreenshot.png] from container [RackCLI]

    Finished! Deleted 3 objects in 1 second

``set-metadata``
~~~~~~~~~~~~~~~~
Sets metadata for the specified container.

::

    rack files container set-metadata --name <containerName> --metadata <key1=val1,key2=val2,...> [optional flags]

**Response**

.. code::

    $ rack files container set-metadata --name RackCLI --metadata heat=true
    Metadata:Heat	true

``get-metadata``
~~~~~~~~~~~~~~~~
Retrieves the metadata for a specified container.

::

    rack files container get-metadata --name <containerName> [optional flags]

**Response**

.. code::

    $ rack files container get-metadata --name RackCLI
    Metadata:Heat	true

``update-metadata``
~~~~~~~~~~~~~~~~~~~
Updates metadata items for a specified container, or adds the specified metadata if no metadata is currently associated with the container.

::

    rack files container update-metadata --name <containerName> --metadata <key1=val1,key2=val2,...> [optional flags]

**Response**

.. code::

    $ rack files container update-metadata --name RackCLI --metadata heat=false
    Metadata:Heat	false

``delete-metadata``
~~~~~~~~~~~~~~~~~~~
Deletes one or more metadata keys for a container.

::

    rack files container delete-metadata --name <containerName> --metadata-keys <key1,key2,...> [optional flags]

**Response**

.. code::

    $ rack files container delete-metadata --name RackCLI --metadata-keys heat
    Successfully deleted metadata with keys [Heat] from container [RackCLI].

Object
------

The ``object`` subcommand provides information about and performs actions on the objects in Cloud Files. The ``object`` subcommand uses the following syntax::

    rack files object <action> [optional flags]

The following sections describe the actions that you can perform on the ``object`` subcommand and provide example responses.

``list``
~~~~~~~~
Lists all of the objects contained in a specified container.

::

    rack files object list --container <containerName> [optional flags]
    (echo containerName1 && echo containerName2) | rack files object list --stdin container [optional flags]

**Response**

.. code::

    $ rack files object list --container RackCLI
    Name			            Bytes	ContentType	LastModified
    arrowsscreenshot.png	74288	image/png	  2015-08-10T17:22:04.149420
    colors-logos.png	    18841	image/png	  2015-08-10T17:22:04.205360
    dashed-lines.png	    54014	image/png	  2015-08-10T17:22:04.165600

``upload``
~~~~~~~~~~
Uploads an object into a specified container.

::

    rack files object upload --container <containerName> --name <objectName> --content <contentToUpload> [optional flags]
    rack files object upload --container <containerName> --name <objectName> --file <fileToUpload> [optional flags]
    cat fileToUpload.txt | rack files object upload --container <containerName> --name <objectName> --stdin content [optional flags]

**Response**

.. code::

    $ rack files object upload --container RackCLI --name Image --file /Users/nath8916/Downloads/rackspace_diagram_icons/PNGs/dedicated-device-generic-1.png
    Successfully uploaded object [Image] to container [RackCLI]

``upload-dir``
~~~~~~~~~~~~~~
Uploads an object directory into a specified container.

::

    rack files object upload-dir --container <containerName> --dir <dirPath> [optional flags]
    find . -type d -name "something*" | rack files object upload-dir --container <containerName> --stdin dir [optional flags]

**Response**

.. code::

    $ rack files object upload-dir --container RackCLI --dir /Users/nath8916/Downloads/rackspace_diagram_icons/PNGs
    Uploaded calendar-small.png to RackCLI
    Uploaded dedicated-server-large.png to RackCLI
    Uploaded file-small.png to RackCLI
    Uploaded dedicated-server-small.png to RackCLI
    Uploaded dedicated-big-data.png to RackCLI
    Uploaded dedicated-device-generic-2.png to RackCLI
    Uploaded dedicated-device-generic-3.png to RackCLI
    Uploaded files-large.png to RackCLI
    Uploaded dedicated-device-generic-1.png to RackCLI
    ...
    Finished! Uploaded 152 objects totaling 4.4MB in 1 minute

``download``
~~~~~~~~~~~~
Downloads an object from the specified container to your local system.

::

    rack files object download --container <containerName> --name <objectName> [optional flags]

``get``
~~~~~~~
Retrieves information about an object.

::

    rack files object get --container <containerName> --name <objectName> [optional flags]

**Response**

.. code::

    $ rack files object get --container RackCLI --name Image
    Name			           Image
    ContentDisposition
    ContentEncoding
    ContentLength		     17763
    ContentType		       application/octet-stream
    StaticLargeObject	   false
    ObjectManifest
    TransID			         tx607965cf09ce42c0b6c13-0055c8e2d6dfw1

``delete``
~~~~~~~~~~
Permanently deletes an object.

::

    rack files object delete --container <containerName> --name <objectName> [optional flags]
    (echo objectName1 && echo objectName2) | rack files object delete --container <containerName> --stdin name [optional flags]

**Response**

.. code::

    $ rack files object delete --container RackCLI --name Image
    Successfully deleted object [Image] from container [RackCLI]

``set-metadata``
~~~~~~~~~~~~~~~~
Sets metadata for the specified object.

::

    rack files object set-metadata --container <containerName> --name <objectName> --metadata <key1=val1,key2=val2,...> [optional flags]

**Response**

.. code::

    $ rack files object set-metadata --container RackCLI --name Image --metadata heat=true
    Metadata:Heat	true

``get-metadata``
~~~~~~~~~~~~~~~~
Retrieves the metadata for a specified object.

::

    rack files object get-metadata --container <containerName> --name <objectName> [optional flags]

**Response**

.. code::

    $ rack files object get-metadata --container RackCLI --name Image
    Metadata:Heat	true

``update-metadata``
~~~~~~~~~~~~~~~~~~~
Updates metadata items for a specified object, or adds the specified metadata if no metadata is currently associated with the object.

::

    rack files object update-metadata --container <containerName> --name <objectName> --metadata <key1=val1,key2=val2,...> [optional flags]

**Response**

.. code::

    $ rack files object update-metadata --container RackCLI --name Image --metadata heat=false
    Metadata:Heat	false

``delete-metadata``
~~~~~~~~~~~~~~~~~~~
Deletes one or more metadata keys from an object.

::

    rack files object delete-metadata --container <containerName> --name <objectName> --metadata-keys <key1,key2,...> [optional flags]

**Response**

.. code::

    $ rack files object delete-metadata --container RackCLI --name Image --metadata-keys heat
    Successfully deleted metadata with keys [Heat] from object [Image].

Account
-------

The ``account`` subcommand provides information about and performs actions on the metadata for your Cloud Files account. The ``account`` subcommand uses the following syntax::

    rack files account <action> [optional flags]

The following sections describe the actions that you can perform on the ``account`` subcommand and provide example responses.

``set-metadata``
~~~~~~~~~~~~~~~~
Sets metadata for the account.

::

    rack files account set-metadata --metadata <key1=val1,key2=val2,...> [optional flags]

**Response**

.. code::

    $ rack files account set-metadata --metadata Temp-Url-Key=asdf1234
    Metadata:Temp-Url-Key	asdf1234

``get-metadata``
~~~~~~~~~~~~~~~~
Retrieves the metadata for the account.

::

    rack files account get-metadata [optional flags]

**Response**

.. code::

    $ rack files account get-metadata
    Metadata:Temp-Url-Key	asdf1234

``update-metadata``
~~~~~~~~~~~~~~~~~~~
Updates metadata items for the account, or adds the specified metadata if no metadata is currently associated with the account.

::

    rack files account update-metadata --metadata <key1=val1,key2=val2,...> [optional flags]

**Response**

.. code::

    $ rack files account update-metadata --metadata Temp-Url-Key=asdf12345
    Metadata:Temp-Url-Key	asdf12345

``delete-metadata``
~~~~~~~~~~~~~~~~~~~
Deletes one or more metadata keys from the account.

::

    rack files account delete-metadata  --metadata-keys <key1,key2,...> [optional flags]

**Response**

.. code::

    $ rack files account delete-metadata --metadata-keys Temp-Url-Key
    Successfully deleted metadata with keys [Temp-Url-Key] from account.

Large-object
------------

Large objects are files larger than 5 GB. Given the designated size of each piece, ``rack`` divides the file into the required number of pieces, appropriately names them, and uploads them to the specified container. Downloading a large object is done with the regular ``rack files object download`` command. 

.. note::

    Although files larger than 5 GB must use the ``large-object`` subcommand, files smaller than 5 GB can also use it.

The ``large-object`` subcommand uses the following syntax::

    rack files large-object <action> [optional flags]
    
The following sections describe the actions that you can perform on the ``large-object`` subcommand and provide example responses.

``upload``
~~~~~~~~~~
Upload a large object to a specified container. Use the ``--size-pieces`` flag to specify the size of the pieces (in MB) to divide the file into. 

::

    rack files large-object upload --container <containerName> --size-pieces <sizePieces> [--name <objectName> | --stdin file] [optional flags]

**Response**

.. code::

    $ rack files large-object upload --container RackCLI --name largeObject --file largeZipFile.zip --size-pieces 500
    Finished! Uploaded object [largeObject] to container [RackCLI] in 5 minutes

``delete``
~~~~~~~~~~
Deletes a large object from a specified container.

::

    rack files large-object delete --container <containerName> --name <objectName>
    (echo objectName1 && echo objectName2) | rack files large-object delete --container <containerName> --stdin name

**Response**

.. code::

    $ rack files large-object delete --container RackCLI --object largeObject
    Deleted object [largeObject] from container [RackCLI]
    
