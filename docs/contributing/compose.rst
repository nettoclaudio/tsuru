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

Running Docker Compose
----------------------

First get the up-to-date Tsuru's source code available on GitHub and enter into the newly created directory.

.. code:: bash

   $ git clone https://github.com/tsuru/tsuru.git
   $ cd tsuru

Run the Docker Compose to up the Tsuru API and its required services. On the first time, this action may take a long time, be patient.
 
.. code:: bash

   $ docker-compose up -d

.. NOTE::
    Everytime you change the Tsuru API and want testing it, you should rebuild that.
    This can be archived by running the command ``docker-compose up --build -d api``.

If everything works as expected, now you have all the needed services running on your machine. You can verify they are running using the command below:

.. code:: bash

   $ docker-compose ps

Every service on this output should be running - with State field `Up`.
For ensuring the Tsuru API is working properly, you can do an HTTP request on its health check endpoint, as shown below:

.. code:: bash

   $ curl http://127.0.0.1:8080/healthcheck
   WORKING

The response message should be "WORKING" meaning the installation was successful.

Creating admin user
-------------------

Now you have a fresh instance of Tsuru installed, you need creating a super-user.
It will be used later to perform any actions on Tsuru.

For creating the administrator user, you must execute an especial command *inside* the Tsuru API container, after hit that it will prompt the password for you. Do that following the command below:

.. code:: bash

    $ docker-compose exec api tsurud root-user-create admin@example.com

Great. You had created the ``admin@example.com`` administrator user. Make sure to remember its email and password.

Now you can be use the newly created Tsuru API into your client, for this creat:

.. code:: bash

    $ tsuru target-add -s development http://127.0.0.1:8080
    $ tsuru login

The login command will prompt the email and password for you, fill that with the super-user credentials created above.

Adding Tsuru Node
------------------

Right now you already are logged in into the local Tsuru API, you need to create one pool of nodes and register the ``node1`` as a Tsuru node.

.. code:: bash

    $ tsuru pool-add development -p -d
    $ tsuru node-add --register address=http://node1:2375 pool=development

You can verify if it was registered properly using the command below:

.. code:: bash

    $ tsuru node-info http://node1:2375

The status of this Tsuru node should be marked as ``ready``.

At this point you are ready to create and deploy apps using Tsuru.

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

