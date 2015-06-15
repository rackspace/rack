.. _home:

Rackspace Command Line Interface
================================

.. warning:: The ``rack`` tool is under heavy development; the name of the binary
             and other portions may rapidly change. If you want to participate or
             provide early feedback, see the `Github project`_

Description
-----------

The Rackspace Command Line Interface is a unified tool to manage your Rackspace
services. It provides streamlined and secure configuration as well as a single
point of entry for all Rackspace Cloud services.

See :ref:`installation_and_configuration` to get started.

Synopsis
--------

::

  rack [--options] <service> <command> <subcommand> [--flags]

All ``rack`` commands follow the pattern above - for example, if you wanted to
list all running servers on your Rackspace account, you would type::

  rack servers instance list

And the response (**default**: table-based output) would look like::

      ID	Name		Status	Public IPv4	Private IPv4	Image	Flavor
      GUID	my_server	ACTIVE	101.130.19.31	10.208.128.233	GUID	io1-30

Options
-------

The ``rack`` CLI uses global options (``[--options]``) to alter the output, or
pass in other required **global** information into the tool, these are:

``--json``
  (boolean) Return output in JSON format.

``--table``
  Return output in tabular format. *This is the default output format.*

``--csv``
  Return output in csv format.

``--help, -h``
  Show help

``--generate-bash-completion``
  Generate bash completion directives for tab-completion of commands.

``--version, -v``
  Print the version

For more on options, see :ref:`global_options`.

Services
--------

* :ref:`servers` - Commands for Rackspace Cloud Servers, dedicated and virtual.


.. toctree::
   :caption: Table of Contents
   :name: mastertoc
   :maxdepth: 2

   self
   configuration.rst
   globaloptions.rst
   services/index.rst



Indices and tables
==================

* :ref:`genindex`
* :ref:`search`

.. _Github project: https://github.com/jrperritt/rack
