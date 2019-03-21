.. Copyright 2017 tsuru authors. All rights reserved.
   Use of this source code is governed by a BSD-style
   license that can be found in the LICENSE file.

++++++++++++++++++++++++++++++++++++++++++++++++++++++
Building a development environment with Docker Compose
++++++++++++++++++++++++++++++++++++++++++++++++++++++

This guide shows how to run Tsuru on a single host using Docker Compose. That can be useful for development and test environments because it allows quickly to start and down the entire stack of services and requires few hardware resources.

.. WARNING::
   Do not run this installation method on production or any other environments where the Tsuru API should be readily available for a user base - only the developer.

To be able to follow this guide, you need installing the Docker_ (v1.13.0 or later), `Docker Compose`_ (v1.10.0 or later) and the `Tsuru client`_. After getting these tools, make sure they are running correctly on your system.

Installing
----------

First get the up-to-date Tsuru's source code available on GitHub and enter into the newly created directory.

.. code:: bash

   $ git clone https://github.com/tsuru/tsuru.git
   $ cd tsuru

Run the Docker Compose to up the Tsuru API and its required services. On the first time, this action may take a long time running, be patient.
 
.. code:: bash

   $ docker-compose up -d

If everything works as expected, now you have all the needed services running on your machine. You can verify they are running using the command below:

.. code:: bash

   $ docker-compose ps

A similar output is shown below, you can use that as referece to find the correct address of services.

    +--------------+-------+-------------------------+
    | Service      | State | Address                 |
    +==============+=======+=========================+
    | ``api``      | Up    | ``127.0.0.1:8080/TCP``  |
    +--------------+-------+-------------------------+
    | ``mongo``    | Up    | ``127.0.0.1:27017/TCP`` |
    +--------------+-------+-------------------------+
    | ``node1``    | Up    | ``127.0.0.1:2375/TCP``  |
    +--------------+-------+-------------------------+
    | ``planb``    | Up    | ``127.0.0.1:8989/TCP``  |
    +--------------+-------+-------------------------+
    | ``redis``    | Up    | ``127.0.0.1:6379/TCP``  |
    +--------------+-------+-------------------------+
    | ``registry`` | Up    | ``127.0.0.1:5000/TCP``  |
    +--------------+-------+-------------------------+

For ensuring the Tsuru API is working properly, you can do an HTTP request on its health check endpoint, as shown below:

.. code:: bash

   $ curl http://127.0.0.1:8080/healthcheck
   WORKING

The response message should be "WORKING" meaning the installation was successful.

Creating admin user
-------------------

Now you have a fresh instance of Tsuru installed, you need creating the super user.
This user will be used to perform any.

To create the administrator user, use the command below then type your password and confirm it.

.. code:: bash

    $ docker-compose exec api tsurud root-user-create admin@example.com

Great. You have been created the ``admin@example.com`` administrator user.

On your tsuru 


Adding Docker Node
------------------

Then configure the tsuru target:

::

    $ tsuru target-add development http://127.0.0.1:8080 -s

You need to create one pool of nodes and add node1 as a tsuru node.
::

    $ tsuru pool-add development -p -d
    $ tsuru node-add --register address=http://node1:2375 pool=development

Everytime you change tsuru and want to test you need to run ``build-compose.sh`` to build tsurud, generate and run the new api.

If you want to use gandalf, generate one app token and insert into docker-compose.yml file in gandalf environment TSURU_TOKEN.

::

    $ docker-compose stop api
    $ docker-compose run --entrypoint="/bin/sh -c" api "tsurud token"
    // insert token into docker-compose.yml
    $ docker-compose up -d

.. _Docker:  https://docs.docker.com/engine/installation/
.. _`Docker Compose`: https://docs.docker.com/compose/install/
.. _Tsuru: https://github.com/tsuru/tsuru
.. _`Tsuru client`: https://tsuru-client.readthedocs.io/en/latest/installing.html

Kubernetes Integration
----------------------

One can register a minikube instance as a cluster in tsuru to be able to orchestrate tsuru applications on minikube.

Start minikube:

::

    $ minikube start --insecure-registry=10.0.0.0/8

Create a pool in tsuru to be managed by the cluster:

::

    $ tsuru pool add kubepool --provisioner kubernetes


Register your minikube as a tsuru cluster:

::

    $ tsuru cluster add minikube kubernetes --addr https://`minikube ip`:8443 --cacert $HOME/.minikube/ca.crt --clientcert $HOME/.minikube/apiserver.crt --clientkey $HOME/.minikube/apiserver.key --pool kubepool

Check your node IP:

::

    $ tsuru node list -f tsuru.io/cluster=minikube

Add this IP address as a member of kubepool:

::

    $ tsuru node update <node ip> pool=kubepool

You are ready to create and deploy apps kubernetes.

