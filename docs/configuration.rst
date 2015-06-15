.. _installation_and_configuration:

Installation and Configuration
==============================

.. warning:: The installation / configuration instructions here are to be considered
             buyer-beware. As development continues a range of configuration options
             will be supported; as such what works today may not work tomorrow, but
             it works right now so that's ok.

Installation
------------

The Rackspace CLI ``rack`` is a self-contained binary writted in go_ this means
that installation is as simple as downloading the relevant binary for your operating
system and ensuring it is on your path.

Binary Downloads:

* `Mac OSX (64 bit)`_
* `Linux (64 bit)`_
* `Windows (64 bit)`_

OSX and Linux
^^^^^^^^^^^^^

After downloading the binary on OSX and Linux you will need to make the binary
executable by typing::

    chmod a+x /path/to/rack

We also recommend you move or symbolically link it on these platforms to `/usr/local/bin`::

    mkdir -p /usr/local/bin/
    ln -s /path/to/rack /usr/local/bin/rack

You can now add it to your path with::

    export PATH=$PATH:/usr/local/bin

Windows
^^^^^^^

After downloading the binary on Windows, you can immediately run it.

We recommend that you copy it to a location outside of your Downloads folder (e.g. C:\\tools) and add that location to your PATH. You must open a new command prompt after modifying the PATH variable.

1. Create a new directory for command line tools, if you don't already have one, e.g. C:\\tools.
2. Copy rack.exe to that directory 
3. Add the directory to your user's PATH environment variable, e.g. ``setx path "%path%;C:\tools"`` or press the Windows key, type "set env", select "Edit environment variables for your account", select the PATH user variable and append ";C:\\tools" to the value and save your changes.
4. Open a new command prompt after modifying the PATH variable.


Configuration
-------------

.. warning:: This will change. I warned you.

In order for the ``rack`` tool to pick up your configuration, you will need to
export the following environment variables (case matters):

* RS_REGION_NAME (DFW, IAD, ORD, LON, SYD, HKG)
* RS_USERNAME (Your Rackspace username)
* RS_API_KEY (Your Rackspace API key)

So, for example on OSX and Linux; you would type::

    export RS_REGION_NAME=IAD
    export RS_USERNAME=<your rackspace username>
    export RS_API_KEY=<secrets>

On Windows you would type::

    set RS_REGION_NAME=IAD
    set RS_USERNAME=<your rackspace username>
    set RS_API_KEY=<secrets>

You can get your API key by logging into the `Cloud Control panel`_ and clicking
on *account -> account settings* and clicking "show" next to "API Key". Be careful;
this key is special. Don't share it! [#]_ [#]_

Advanced Configuration Values
-----------------------------

Identity Endpoint
^^^^^^^^^^^^^^^^^

If you require pointing to a custom Cloud Identity endpoint; you can set the
following environment variable:

* RS_AUTH_URL (https://identity.api.rackspacecloud.com/v2.0)

For example::

    export RS_AUTH_URL=https://identity.api.rackspacecloud.com/v2.0

.. [#] **No, seriously** - don't share it, don't check it into source control, the API
      gives anyone who has it god-like powers. We've accidentally shared it in the
      past and literally had a rip in space-time that sucked us into an alternate
      dimension that has mutant pug overlords. Please save us.

.. [#] Hush now human. No tears. Only sleep.

.. _go: https://golang.org/
.. _Mac OSX (64 bit): https://ba7db30ac3f206168dbb-7f12cbe7f0a328a153fa25953cbec5f2.ssl.cf5.rackcdn.com/Darwin/amd64/rack
.. _Linux (64 bit): https://ba7db30ac3f206168dbb-7f12cbe7f0a328a153fa25953cbec5f2.ssl.cf5.rackcdn.com/Linux/amd64/rack
.. _Windows (64 bit): https://ba7db30ac3f206168dbb-7f12cbe7f0a328a153fa25953cbec5f2.ssl.cf5.rackcdn.com/Windows/amd64/rack.exe
.. _Cloud Control panel: https://mycloud.rackspace.com/
