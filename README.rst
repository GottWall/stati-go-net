Welcome to stati-go-net's documentation!
========================================

stati-go-net is a golang client with HTTP & UDP & TCP/IP  transports for `GottWall metrics aggregation platform <http://github.com/GottWall/GottWall>`_

.. image:: https://secure.travis-ci.org/GottWall/stati-go-net.png
	   :target: https://secure.travis-ci.org/GottWall/stati-go-net


INSTALLATION
------------

To use stati_go execute next command::

  go install github.com/GottWall/stati_go_net


It's install ``stati_net`` package to ``GOROOT`` path.



USAGE
-----

To use library in your packages import it first.

.. sourcecode:: go

   package your_package_name

   import (
      "github.com/GottWall/stati-go-net"
   )

   var (
       solt_base int    = DEFAULT_SOLT_BASE
       project string = "test_project"
       private_key string = "private_key"
       public_key  string = "public_key"
       host string = "127.0.0.1"
       port int16 = 8890
       proto string = "https"
       prefix string = "/custom_prefix"
	)


HTTPClient
^^^^^^^^^^

.. sourcecode:: go

   // HTTPClient
   // stati_net package already imported
   func main() {

   var client *stati_net.HTTPClient = stati_net.HTTPClientInit(
      project, private_key, public_key, host, port, proto, prefix)

   client.Incr("name", 10, 1392454739)
   }


TCP/IP Client
^^^^^^^^^^^^^

.. sourcecode:: go

   // TCP/IP Client
   // stati_net package already imported

   func main() {
      var client *stati_net.TCPClient = stati_net.TCPClientInit(
         project, private_key, public_key, host, port,
	     auth_delimiter, chunk_delimiter)

      client.Incr("name", 10, 1392454739)
   }



UDP Client
^^^^^^^^^^

.. sourcecode:: go

   // UDP client
   // stati_net package already imported

   func main() {
      var client *stati_net.UDPClient = stati_net.UDPClientInit(
         project, private_key, public_key, host, port,
		 auth_delimiter, chunk_delimiter)

	  // Change default UDP packet size
	  client.SetMaxPacketSize(4096)
      client.Incr("name", 10, 1392454739)
   }



CONTRIBUTE
----------

We need you help.

#. Check for open issues or open a fresh issue to start a discussion around a feature idea or a bug.
   There is a Contributor Friendly tag for issues that should be ideal for people who are not very familiar with the codebase yet.
#. Fork `the repository`_ on Github to start making your changes to the **develop** branch (or branch off of it).
#. Write a test which shows that the bug was fixed or that the feature works as expected.
#. Send a pull request and bug the maintainer until it gets merged and published.

.. _`the repository`: https://github.com/GottWall/stati-go-net/
