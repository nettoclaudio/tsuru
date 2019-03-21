.. Copyright 2017 tsuru authors. All rights reserved.
   Use of this source code is governed by a BSD-style
   license that can be found in the LICENSE file.

++++++++++++++++++++++++++++++++++++++++++++++++++++++
Building a development environment with Docker Compose
++++++++++++++++++++++++++++++++++++++++++++++++++++++

This guide shows how to run Tsuru on a single host using Docker Compose.
This installation method can be useful for development and test environments where the server does not need to be readily available for any user - only the developer.

.. WARNING::
    Do not run this installation method on production environments.

To be able to follow this guide, you need installing the Docker_ (v1.13.0 or later) and `Docker Compose`_ (v1.10.0 or later). After getting these tools, make sure they are running correctly on your system.

Running Docker Compose
----------------------

Get the up-to-date Tsuru's source code available on GitHub and enter into that directory.

.. code:: bash
    $ git clone https://github.com/tsuru/tsuru.git
    $ cd tsuru

Then run the Docker Compose to up the Tsuru API and its required services. At first time this action may take a long time running, be patient.
 
.. code:: bash
    $ docker-compose up -d

If everything works as expected, now you have Tsuru dependencies (such as MongoDB and Redis databases, PlanB application router and Registry), the Tsuru API and one Docker Node all them running in your machine. You can verify they are running using the command below:

.. code:: bash
    $ docker-compose ps


You have a fresh tsuru installed, so you need to create the admin user running tsurud inside container.

::

    $ docker-compose exec api tsurud root-user-create admin@example.com

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

