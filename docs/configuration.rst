.. _installation_and_configuration:

Installation and configuration
==============================

This section provides complete and detailed information for installing and configuring the Rackspace Command Line Interface (``rack`` CLI). 

Install the CLI
---------------

The ``rack`` CLI is a self-contained binary written in go_. To install the CLI, you simply download the relevant binary for your OS and ensure that the directory in which it resides is in your system’s PATH environment variable.

Download the binary for your OS from one of the following links, and then follow the specific instructions for your OS:

* `Mac OS X (64-bit)`_
* `Linux (64-bit)`_
* `Windows (64-bit)`_

Mac OS X and Linux with Homebrew
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

If you are using `Homebrew`_, run the following command::

    brew install rack

Mac OS X and Linux without Homebrew
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

1. After downloading the binary for Mac OS X or Linux, make the binary executable by running the following command::

    chmod a+x /path/to/rack

2. Move it to or symbolically link it to ``/usr/local/bin``, as follows::

    mkdir -p /usr/local/bin/
    ln -s /path/to/rack /usr/local/bin/rack

3. Add it to your PATH environment variable as follows::

    export PATH=$PATH:/usr/local/bin

Windows with Chocolatey
^^^^^^^^^^^^^^^^^^^^^^^

If you are using `Chocolatey`_, run the following command::

    choco install rack

Windows without Chocolatey
^^^^^^^^^^^^^^^^^^^^^^^^^^

You can install the binary manually or by using a script.

Install the binary manually
~~~~~~~~~~~~~~~~~~~~~~~~~~~

After you download the binary on Windows, you can immediately run it.

We recommend that you copy the binary to a location outside of your **Downloads** folder and add that location to your PATH environment variable, as follows:

1. If you don’t already have one, create a new directory for command-line tools. For example, ``C:\tools``.
2. Copy **rack.exe** to that directory.
3. Add the directory to your user's PATH environment variable. You can do this by opening a command prompt window and run the following command::
    
    setx path "%path%;C:\tools"
    
4. After modifying the PATH variable, open a new command prompt window.

Install the binary with a script
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

This option requires PowerShell version 3 or later.

Perform the following steps to set up the ``rack`` CLI. Alternatively, you can save the script as a PowerShell file (for example, **rackspace-cli.ps1**) and execute it.

Open PowerShell_Ise, paste the following script in the scripting pane, and then click the green play button to start the execution.

::

    #requires -Version 3
    $DownloadPath = 'C:\Tools'
    
    Write-Output -InputObject "[$(Get-Date)] Status  :: Set the Tools Directory $DownloadPath"
    New-Item -Path $DownloadPath -ItemType Directory -ErrorAction SilentlyContinue > $null
    Set-Location -Path $DownloadPath -ErrorAction SilentlyContinue
    
    Write-Output -InputObject "[$(Get-Date)] Status  :: Download Rackspace CLI in C:\Tools"
    Invoke-WebRequest -Uri 'https://goo.gl/NMvmcx/Windows/amd64/rack.exe' -Method Get -OutFile rack.exe
     
    Write-Output -InputObject "[$(Get-Date)] Status  :: Unblock the executable file rack.exe"
    Unblock-File -Path $("$DownloadPath\rack.exe")
    
    Write-Output -InputObject "[$(Get-Date)] Status  :: Permanently set the path $DownloadPath to the Environment variable (Reboot required)."
    [System.Environment]::SetEnvironmentVariable('Path', $env:Path + 'C:\Tools', [System.EnvironmentVariableTarget]::Machine)
    Write-Output -InputObject "[$(Get-Date)] Status  :: Temporarily set the path $DownloadPath to the Environment variable for immediate use in the current powershell session"
     $env:Path += ';C:\Tools'

Configuration
-------------

To authenticate against the Rackspace API, the following parameters are required:

* Username: your Rackspace username
* API key: your Rackspace API key
* Region: a Rackspace region
* Authentication Endpoint (URL): (Optional/Advanced) A URL to send the authentication request to

You can specify these parameters quickly by using the interactive ``configure`` command (recommended), or you can use other methods such as specifying command-line flags, manually creating or editing a configuration file, or setting environment variables. All of these methods are explained in this section.

Interactive configure command
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

If this is your first time using the ``rack`` CLI, we recommend that you run the interactive ``configure`` command.

.. note::
   Windows users should use PowerShell, not PowerShell ISE to run this command.

The ``configure`` command automatically creates a configuration file for you if one doesn't already exist and walks you through creating a profile for it::

    rack configure

    This interactive session will walk you through creating
    a profile in your configuration file. You may fill in all or none of the
    values.

    Rackspace Username: <yourRackspaceUsername>
    Rackspace API key: <yourRackspaceApiKey>
    Rackspace Region: <theRackspaceRegion>
    Profile Name (leave blank to create a default profile):
    
Username is the username for your Rackspace Cloud account. You can get your API key by logging in to the Cloud Control Panel, clicking on your account name in the upper-right corner, and then selecting **Account Settings**. The region is the region where your Rackspace infrastructure is deployed. If you want to create a profile other than the default profile, enter a name for the profile.

After the profile is created, you can immediately start working. For example, you could issue the following command to get a list of the servers on your Rackspace account::

    rack servers instance list

Command-line options
^^^^^^^^^^^^^^^^^^^^

If used, the following command-line authentication flags take precedence over any
other parameters located in any other forms of authentication (configuration file and
environment variables).

* ``--username``
* ``--api-key``
* ``--region``

Configuration file
^^^^^^^^^^^^^^^^^^

Any authentication parameters not set on the command line are looked for in a configuration file. The configuration file should be located in ``$HOME/.rack/config``. When you use the interactive ``configure`` command, a configuration file is automatically created. 

The configuration file format is similar to the following format::

    username=<yourRackspaceUsername>
    api-key=<yourRackspaceApiKey>
    region=<theRackspaceRegion>

    [another-profile]
    username=<anotherRackspaceUsername>
    api-key=<anotherRackspaceApiKey>

The preceding example shows a default profile that doesn't have a named section. ``another-profile`` is a different profile in the configuration file. When you use the default profile, you don't need to supply a flag when executing ``rack``. You can specify a profile on the command line with the ``profile`` flag.

::

    rack servers instance list --profile another-profile

Note that none of the authentication parameters have to be set in the configuration file. Parameters not set there are looked for elsewhere.

Environment variables
^^^^^^^^^^^^^^^^^^^^^

Finally, ``rack`` looks for any remaining unset authentication parameters in environment variables. The following are values are permitted (case matters):

* ``RS_REGION_NAME``: DFW, IAD, ORD, LON, SYD, HKG
* ``RS_USERNAME``: your Rackspace username
* ``RS_API_KEY``: your Rackspace API key

For example, on OS X and Linux, you would type::

    export RS_REGION_NAME=IAD
    export RS_USERNAME=yourRackspaceUsername
    export RS_API_KEY=yourRackspaceApiKey

On Windows, you would type::

    set RS_REGION_NAME=IAD
    set RS_USERNAME=yourRackspaceUsername
    set RS_API_KEY=yourRackspaceApiKey

Command completion
------------------
To set up command completion for the Bash shell, run ``rack init``.

Currently, this command is available only for the Bash shell. If you're using a Linux OS, ``rack init`` will look for and, if found, amend ``$HOME/.bashrc`` to enable command completion. If you're on a Darwin OS (like Mac), it will look for ``$HOME/.bash_profile``.

If you want to set up command completion yourself (or if you're on a Windows OS and using a Bash shell), you can copy the following file to the appropriate directory and source it:
`https://github.com/rackspace/rack/blob/master/setup/commandcompletion_bash.sh`

If you are using PowerShell and want command completion, you can run the ``commandcompletion_posh.ps1`` script, also located in the ``setup`` directory. That script performs normal command completion for non-``rack`` commands, and completions for ``rack`` commands. A few caveats for PowerShell users:

* The script overrides the ``global:TabExpansion2`` function.
* The script should work for PowerShell versions later than or equal to 3, but it was tested with PowerShell_ISE v4.
* You get the normal Windows command completion (with a circular buffer).

Check the version
-----------------

To see the current version of the CLI, run the following command::

    rack version

    rack version 0.0.0-dev
    commit: d69f4d2030b307076ad0a10f4b5addf440493aec

Advanced configuration values
-----------------------------

If you need to point to a custom Cloud Identity endpoint, you can set the following environment variable::

    RS_AUTH_URL=https://identity.api.rackspacecloud.com/v2.0

For example::

    export RS_AUTH_URL=https://identity.api.rackspacecloud.com/v2.0

In addition, you can provide it as a flag on the command-line or as a value in the configuration file profile. In either case, the parameter name is ``auth-url``.




.. _go: https://golang.org/
.. _Mac OS X (64-bit): https://ec4a542dbf90c03b9f75-b342aba65414ad802720b41e8159cf45.ssl.cf5.rackcdn.com/1.1.1/Darwin/amd64/rack
.. _Linux (64-bit): https://ec4a542dbf90c03b9f75-b342aba65414ad802720b41e8159cf45.ssl.cf5.rackcdn.com/1.1.1/Linux/amd64/rack
.. _Windows (64-bit): https://ec4a542dbf90c03b9f75-b342aba65414ad802720b41e8159cf45.ssl.cf5.rackcdn.com/1.1.1/Windows/amd64/rack.exe
.. _Homebrew: http://brew.sh
.. _Chocolatey: http://chocolatey.org
.. _Cloud Control panel: https://mycloud.rackspace.com/
