# Apocryph 
Consensus Network for Autonomous Agents

[![Discord](https://img.shields.io/badge/DISCORD-COMMUNITY-informational?style=for-the-badge&logo=discord)](https://discord.gg/ESr9KMR)

> Apocryph Agents can automate the cash flow in autonomous organizations, optimize city traffic, or reward the computing power used to train their own neural networks.
## Table of Contents

- [Overview](#overview)
  - [Quick Summary](#quick-summary)
- [Getting Started](#getting-started)
  - [Prerequisite](#prerequisite)
  - [Create agent project](#create-agent-project)
  - [Enable agent testbed](#enable-agent-testbed)
  - [Configure agent testbed](#configure-agent-testbed)
  - [Create your agents](#create-your-agents)
  - [Run your first multi-agent distributed application](#run-your-first-multi-agent-distributed-application)
- [How Apocryph works](#how-apocryph-works)
  - [Agent Model](#agent-model)
    - [Reminders](#reminders)
    - [Publish and Subscribe](#publish-and-subscribe)
    - [Object Capability Security Model](#object-capability-security-model)
    - [Call Tickets](#call-tickets)
    - [Messages](#messages)
    - [Services](#services)
  - [Consenus](#consenus)
    - [Inter Blockchain Communication](#inter-blockchain-communication)
  - [Network Nodes](#network-nodes)
    - [Scalability](#scalability)
- [Apocryph Developer Node](#apocryph-developer-node)
  - [Running via Docker Compose](#running-via-docker-compose)
  - [Running natively](#running-natively)

## Overview

Apocryph is a new consensus network for autonomous agents. From developer perspective,
we have put a great focus on selecting a technology stack comprising widely adopted platforms,
tools and development paradigms.

Below, you can see a short video of how easy it is to setup Apocryph test node on your 
local development machine using only Docker and Docker-Compose:

[![asciicast](https://asciinema.org/a/295036.svg)](https://asciinema.org/a/295036?speed=2)

### Quick Summary

Apocryph is an architecture:

- defines patterns and practices for building distributed systems
- covers both open-source and closed-source parts of the system being built
- compliant with the latest enterprise-grade software architectures and technologies

Apocryph is a framework:

- has built-in library for building multi-agent systems
- supports both proactive and passive agents

Apocryph is a blockchain *(implementation in-progress)*:

- implements highly scalable DPoS BFT consensus 
- designed to be inter-blockchain communication ready

Apocryph is an economy *(implementation in-progress)*:

- supports fully programmable digital economy model
- accommodates both humans and AI actors 

## Getting Started

### Prerequisite

### Create agent project

### Enable agent testbed

### Configure agent testbed

### Create your agents

### Run your first multi-agent distributed application

## How Apocryph works

![Architecture Overview](docs/images/architecture_overview.jpg "Architecture Overview")
 
### Agent Model

#### Reminders

#### Publish and Subscribe

#### Object Capability Security Model

#### Call Tickets

#### Messages

#### Services

Apocryph services are comprised of custom logic which allows Agents to communicate with 
the outside world. They allow one to extend the consensus algorithm and provide additional 
ways to receive inputs and produce outputs for an agent.

While service execution is not covered by consensus between nodes (in the way agent execution is), 
the different instances of services running on different nodes are expected to give the same outputs 
on most nodes, so that the nodes can reach consensus on the input to the agent.

### Consenus

Apocryph consensus implementation is using serverless, stream-based architecture to 
achieve high concurrency and throughput. For intra-node communication it is using [Peprer](https://github.com/obecto/perper) 
and for inter-node communication and persistence it is using [IPFS](https://ipfs.io/).

#### Inter Blockchain Communication

### Network Nodes

Apocryph is built on top of [Peprer](https://github.com/obecto/perper) - stream-based, horizontally 
scalable framework for asynchronous data processing. This allows Apocryph Nodes to both
work on a single machine (using docker-compose) or in a datacenter grade cluster environment
using [Kubernetes](http://kubernetes.io/).

#### Scalability

## Apocryph Developer Node

You can test and experiment with Apocryph runtime by running a local
developer node. You can run this node on all major operating systems: 
Windows, Linux and macOS.

### Running via Docker Compose
Using Docker Compose to run Apocryph runtime is the recommended way for users that
would like to run Apocryph Developer Node.

### Prerequisite
- Install [Docker](https://docs.docker.com/install/)
- Install [Docker Compose](https://docs.docker.com/compose/install/)

#### Start IPFS Daemon

Apocryph uses IPFS for its DPoS consensus implementation, thus requires IPFS daemon to run locally on the node:

```bash
docker-compose up -d ipfs
```

#### Start Apocryph Runtime

Before running the Apocryph runtime locally you have to start Perper Fabric in local 
development mode:

- Create Perper Fabric IPC directory  
```bash
mkdir -p /tmp/perper
```
- Run Perper Fabric Docker (This steps require pre-built Perper Fabric image. More information can be found [here](https://github.com/obecto/perper))
```bash
docker-compose up -d perper-fabric
```

Apocryph runtime is implemented as Azure Functions App and can be started with:
```bash
docker-compose up apocryph-runtime
```

### Running natively

In addition to using Docker Compose, you can run Apocryph Developer Node 
natively on your machine. This setup is recommended if you are doing 
source code contributions to Apocryph Runtime. The recommended operating 
system for this setup is Ubuntu 18.04 LTS. 

#### Prerequisite

Before running this sample, you must have the following:

- Install [Azure Functions Core Tools v3](https://docs.microsoft.com/en-us/azure/azure-functions/functions-run-local#v2)
- Install [.NET Core SDK 3.1](https://dotnet.microsoft.com/download/dotnet-core/3.1)
- Install [Docker](https://docs.docker.com/install/)
- Install [IPFS](https://ipfs.io/#install)

#### Enable Perper Functions

Apocryph is based on [Perper](https://github.com/obecto/perper) - stream-based,
horizontally scalable framework for asynchronous data processing. To run Apocryph 
make sure you have cloned Perper repo and have the correct path in Apocryph.proj file.

#### Start IPFS Daemon

Apocryph uses IPFS for its DPoS consensus implementation, thus requires IPFS daemon to run locally on the node:

```bash
ipfs daemon --enable-pubsub-experiment
```

#### Start Apocryph Runtime

Before running the Apocryph runtime locally you have to start Perper Fabric in local 
development mode:

- Building Perper Fabric Docker (in the directory where Perper repo is cloned)
```bash
docker build -t perper/fabric -f docker/Dockerfile .
```
- Create Perper Fabric IPC directory  
```bash
mkdir -p /tmp/perper
```
- Run Perper Fabric Docker 
```bash
docker run -v /tmp/perper:/tmp/perper --network=host --ipc=host -it perper/fabric
```

Apocryph runtime is implemented as Azure Functions App and can be started with:
```bash
func start
```
