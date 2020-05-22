# Apocryph 
Consensus Network for Autonomous Agents

> Apocryph Agents can automate the cash flow in autonomous organizations, optimize city traffic, or reward the computing power used to train their own neural networks.

[![Discord](https://img.shields.io/badge/DISCORD-COMMUNITY-informational?style=for-the-badge&logo=discord)](https://discord.gg/ESr9KMR)

## Table of Contents

- [Overview](#overview)
  - [Quick Summary](#quick-summary)
- [Getting Started](#getting-started)
  - [Prerequisite](#prerequisite)
  - [Create project](#create-project)
  - [Enable testbed](#enable-testbed)
  - [Configure testbed](#configure-testbed)
  - [Create your agents](#create-your-agents)
  - [Run your first multi-agent distributed application](#run-your-first-multi-agent-distributed-application)
- [Apocryph Architecture Overview](#apocryph-architecture-overview)
  - [Agent Model](#agent-model)
    - [State](#state)
    - [Reminders](#reminders)
    - [Publish and Subscribe](#publish-and-subscribe)
    - [Object Capability Security Model](#object-capability-security-model)
    - [Call Balances](#call-balances)
    - [Invocations](#invocations)
  - [Consensus](#consenus)
    - [Selection](#selection)
    - [Querying](#querying)
    - [Gossiping](#gossiping)
    - [Agent Zero](#agent-zero)
    - [Inter Blockchain Communication](#inter-blockchain-communication)
  - [Network Nodes](#network-nodes)
    - [Client and Services](#client-and-services)
    - [Availability](#availability)
- [Test Harness](#test-harness)
- [Contributing](#contributing)

## Overview

Apocryph is a new consensus network for autonomous agents. From developer perspective,
we have put a great focus on selecting a technology stack comprising widely adopted platforms,
tools and development paradigms.

Below, you can see a short video of how easy it is to setup Apocryph test node on your 
local development machine using only Docker and Docker-Compose:

[![asciicast](docs/images/developer_node_rec.png)](https://asciinema.org/a/295036?speed=2&rows=30)

### Quick Summary

Apocryph is an architecture:

- defines patterns and practices for building distributed systems
- covers both open-source and closed-source parts of the system being built
- compliant with the latest enterprise-grade software architectures and technologies

Apocryph is a framework:

- has built-in library for building multi-agent systems
- supports both active and passive agents

Apocryph is a blockchain *(implementation in-progress)*:

- implements highly scalable leaderless consensus 
- designed in mind with inter-blockchain communication

Apocryph is an economy *(implementation in-progress)*:

- supports fully programmable digital economy model
- accommodates both humans and AI actors 

## Getting Started

This is a quick start guide of how to create a simple multi-agent system
using Apocryph.  

### Prerequisite

Before running this guide, you must have the following:

- Install [Azure Functions Core Tools v3](https://docs.microsoft.com/en-us/azure/azure-functions/functions-run-local#v2)
- Install [.NET Core SDK 3.1](https://dotnet.microsoft.com/download/dotnet-core/3.1)
- Install [Docker](https://docs.docker.com/install/)

### Create project

> **NOTE:** As a best practice, the agents should be developed as a separate Class Library 
that is referenced by the function app project.

Run the following command from the command line to create a function app project 
in the SampleApp folder of the current local directory. For simplicity this 
project will contain both the agents source code and testbed configuration (*see the note above*).

```bash
func init SampleApp
```

When prompted, select a worker runtime - for now only dotnet is fully supported.

After the project is created, use the following command to navigate to the new SampleApp project folder.

```bash
cd SampleApp
````
### Enable testbed

To run your agents on your developer machines you can use the 
Apocryph testbed. To use it, you have to clone Apocryph GitHub repo
and add reference to Apocryph.Testbed. 

There are two more NuGet packages that are required:

- Microsoft.Azure.Functions.Extensions
- Microsoft.NET.Sdk.Functions

After theese configurations, your project file will be similar to this:

```xml
<Project Sdk="Microsoft.NET.Sdk">

    <PropertyGroup>
        <TargetFramework>netcoreapp3.1</TargetFramework>
        <AzureFunctionsVersion>v3</AzureFunctionsVersion>
        <LangVersion>8</LangVersion>
        <Nullable>enable</Nullable>
    </PropertyGroup>

    <ItemGroup>
        <None Update="host.json">
            <CopyToOutputDirectory>PreserveNewest</CopyToOutputDirectory>
        </None>
        <None Update="local.settings.json">
            <CopyToOutputDirectory>PreserveNewest</CopyToOutputDirectory>
            <CopyToPublishDirectory>Never</CopyToPublishDirectory>
        </None>
    </ItemGroup>
    
    <ItemGroup>
      <ProjectReference Include="..\..\Apocryph.Testbed\Apocryph.Testbed.csproj" />
    </ItemGroup>
    
    <ItemGroup>
      <PackageReference Include="Microsoft.Azure.Functions.Extensions" Version="1.0.0" />
      <PackageReference Include="Microsoft.NET.Sdk.Functions" Version="3.0.5" />
    </ItemGroup>

</Project>
```

### Configure testbed

Using the testbed requires adding a small portion of boilerplate code that
will enable a local execution of your agents. Using this you can debug
your agents as regular .NET project.

First, you have to enable the testbed and the logging as services. To do this add
Startup.cs file in the root of your project:

```csharp
using Apocryph.Testbed;
using Microsoft.Azure.Functions.Extensions.DependencyInjection;
using Microsoft.Extensions.DependencyInjection;

[assembly: FunctionsStartup(typeof(SampleApp.Startup))]

namespace SampleApp
{
    public class Startup : FunctionsStartup
    {
        public override void Configure(IFunctionsHostBuilder builder)
        {
            builder.Services.AddLogging();
            builder.Services.AddTransient(typeof(Testbed), typeof(Testbed));
        }
    }
}
```

You also have to enable the logging in the host.json:

```json
{
    "version": "2.0",
    "logging": {
        "logLevel": {
            "SampleApp": "Trace"
        }
    }
}
```

Second, you have to create the testbed functions used as main entrypoints
for setting up the agents execution environment. To do this add App.cs file 
in the root of your project:

```csharp
using System.Collections.Generic;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Testbed;
using Apocryph.Agent;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace SampleApp
{
    public class App
    {
        private readonly Testbed _testbed;

        public App(Testbed testbed)
        {
            _testbed = testbed;
        }

        [FunctionName("Setup")]
        public async Task Setup(
            [PerperStreamTrigger(RunOnStartup = true)] PerperStreamContext context,
            CancellationToken cancellationToken)
        {
            await _testbed.Setup(context, "AgentOne", "Runtime", "Monitor", cancellationToken);
        }

        [FunctionName("Runtime")]
        public async Task Runtime(
            [PerperStreamTrigger] PerperStreamContext context,
            [Perper("agentDelegate")] string agentDelegate,
            [PerperStream("commands")] IAsyncEnumerable<AgentCommands> commands,
            CancellationToken cancellationToken)
        {
            await _testbed.Runtime(context, agentDelegate, commands, cancellationToken);
        }

        [FunctionName("Monitor")]
        public async Task Monitor(
            [PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("commands")] IAsyncEnumerable<AgentCommands> commands,
            CancellationToken cancellationToken)
        {
            await _testbed.Monitor(commands, cancellationToken);
        }
    }
}
```
### Create your agents

In the previous step we have configured the testbed entrypoints, by specify
the name of our root agent ("AgentOne"):

```csharp
[FunctionName("Setup")]
public async Task Setup(
    [PerperStreamTrigger(RunOnStartup = true)] PerperStreamContext context,
    CancellationToken cancellationToken)
{
    await _testbed.Setup(context, "AgentOne", "Runtime", "Monitor", cancellationToken);
}
```

You can use any other name that is more suitable for your multi-agent system domain 
(for example: "Organization", "Template" or other). This name indicates the first 
agent you have to create, serving as entrypoint to your multi-agent system.

In the testbed, every agent is represented by a function configured with a small boilerplate
(we will group our agents in a separate namespace called "Agents"). For simplicity we will
colocate the boilerplate (AgentOneWrapper class) and the actual source code (AgentOne class) 
in a single C# file.

*Agents\AgentOne.cs*
```csharp
using System;
using System.Collections.Generic;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Testbed;
using Apocryph.Agent;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace SampleApp.Agents
{
    public class AgentOne
    {
        public Task<AgentContext> Run(object state, AgentCapability self, object message)
        {
            var context = new AgentContext(state, self);
            if (message is AgentRootInitMessage rootInitMessage)
            {
                var cap = context.IssueCapability(new[] {typeof(PingPongMessage)});
                context.CreateAgent("AgentTwo", "AgentTwo", new PingPongMessage {AgentOne = cap}, null);
            }
            else if(message is PingPongMessage pingPongMessage)
            {
                context.SendMessage(pingPongMessage.AgentTwo, new PingPongMessage
                {
                    AgentOne = pingPongMessage.AgentOne,
                    AgentTwo = pingPongMessage.AgentTwo,
                    Content = "Ping"
                }, null);
            }
            return Task.FromResult(context);
        }
    }

    public class AgentOneWrapper
    {
        private readonly Testbed _testbed;

        public AgentOneWrapper(Testbed testbed)
        {
            _testbed = testbed;
        }

        [FunctionName("AgentOne")]
        public async Task AgentOne(
            [PerperStreamTrigger] PerperStreamContext context,
            [Perper("agentId")] string agentId,
            [Perper("initMessage")] object initMessage,
            [PerperStream("commands")] IAsyncEnumerable<AgentCommands> commands,
            [PerperStream("output")] IAsyncCollector<AgentCommands> output,
            CancellationToken cancellationToken)
        {
            await _testbed.Agent(new AgentOne().Run, agentId, initMessage, commands, output, cancellationToken);
        }
    }
}
```

The root agent is a regular agent with the only specific that it receives
a special init message ("AgentRootInitMessage") by the runtime.

The logic for our sample root agent is to create another agent ("AgentTwo")
and start passing back and forward a simple message ("PingPongMessage"). In a 
similar way we can create the source code of our second agent ("AgentTwo").

*Agents\AgentTwo.cs*
```csharp
using System.Collections.Generic;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Testbed;
using Apocryph.Agent;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace SampleAgents.FunctionApp.Agents
{
    public class AgentTwo
    {
        public Task<AgentContext> Run(object state, AgentCapability self, object message)
        {
            var context = new AgentContext(state, self);
            if(message is PingPongMessage initMessage && initMessage.AgentTwo == null)
            {
                var cap = context.IssueCapability(new[] {typeof(PingPongMessage)});
                context.SendMessage(initMessage.AgentOne, new PingPongMessage
                {
                    AgentOne = initMessage.AgentOne,
                    AgentTwo = cap
                }, null);
            }
            else if(message is PingPongMessage pingPongMessage)
            {
                context.SendMessage(pingPongMessage.AgentOne, new PingPongMessage
                {
                    AgentOne = pingPongMessage.AgentOne,
                    AgentTwo = pingPongMessage.AgentTwo,
                    Content = "Pong"
                }, null);
            }
            return Task.FromResult(context);
        }
    }

    public class AgentTwoWrapper
    {
        private readonly Testbed _testbed;

        public AgentTwoWrapper(Testbed testbed)
        {
            _testbed = testbed;
        }

        [FunctionName("AgentTwo")]
        public async Task AgentTwo(
            [PerperStreamTrigger] PerperStreamContext context,
            [Perper("agentId")] string agentId,
            [Perper("initMessage")] object initMessage,
            [PerperStream("commands")] IAsyncEnumerable<AgentCommands> commands,
            [PerperStream("output")] IAsyncCollector<AgentCommands> output,
            CancellationToken cancellationToken)
        {
            await _testbed.Agent(new AgentTwo().Run, agentId, initMessage, commands, output, cancellationToken);
        }
    }
}
```

### Run your first multi-agent distributed application

To run your application you have to first start the Perper Fabric. 
Both Apocryph Runtime and Testbed are using [Perper](https://github.com/obecto/perper)
which is a stream-based, horizontally scalable framework for asynchronous data processing.

You can run Perper Fabric by executing the following command:

```bash
docker run -p 10800:10800 -p 40400:40400 -it obecto/perper-fabric
```

Then you can run your SampleApp as a regular Azure Functions application 
using the following command (in you "SampleApp" project folder):

```bash
func start --build
```

## Apocryph Architecture Overview

From architecture standpoint Apocryph can be viewed as a framework for
developing multi-agent systems running on a decentralized network. The framework
comprises of three main layers: Agent Model, Consensus and Network Nodes.
 
### Agent Model

Multi-agent systems typically consists of number of agents that interact with 
their environment. Apocryph agents follow the same model, they can observe 
the environment by subscribing to the output of other agents, services and 
based on these observations, the agents can emit own publications or pro-actively
engage with other agents. Therefore, Apocryph support both passive, active and
cognitive agents.

Every interaction between the agents and the environment or between the agents
is represented as a command. This allows agents to be executed asynchronously
in a reproducible way.

> **Determinism:** Apocryph agents are implemented using high level languages - 
C# and Python. It is developer responsibility to write deterministic agents that can 
reach consensus on the decentralized network. There are variety of well known practices
and linters for writing deterministic code, for example: *use deterministic seed for 
pseudo-random number generation; avoid floating point types or use them with extra caution; 
avoid random language features (dictionary iterators, uninitialized memory and etc)*.

#### State

Apocryph agents can be viewed as state machine - receive messages, update their state and 
output new messages (in the form of commands). The internal state of the agent is opaque 
object for the network and its structure is known only to its agent owner. 

#### Reminders

Reminded command allows agents to be activated when a specific deadline (point of time)
has passed. Upon activation the agent receives a message 
specified at the time when the command has been emitted. 

> **Time:** There is no guarantee of the time gap between the time of agent activation 
and the requested deadline by the agent.

#### Publish and Subscribe

Publish and subscribe commands allows the agents to observe and change the environment 
by indirectly exchanging messages. Every agent has a public topic associated with the 
agent identifier where the agent can emit arbitrary information and the subscribing agents
gets activated on new publications. An agent can dynamically subscribe to an arbitrary number
of public topics of other agents.

#### Object Capability Security Model

Apocryph agents can directly interact between each other over [object capabilities](https://wiki.c2.com/?ObjectCapabilityModel). Apocryph
object capability contain whitelisted message types and it is created locally by the agent. Then
the agent can distribute the newly created capability by embedding it in messages. Apocryph agents can store 
capabilities (both their own capabilities and capabilities received from other agents) in their 
state for later use.

Security and unforgeability of all object capabilities used by the agents is implemented as part
of the decentralized network that hosts the agents. All object capabilities are embedded in the 
blocks and the network has to agree on their authenticity.     

#### Call Balances

Direct inter-agent interactions incur costs for the receiving agents. To cover for these costs
every agent has associated call balance and the respective costs are deducted form the balance
one every interaction. Therefore the agent initiating the interaction has to transfer (directly or indirectly)
the necessary funds in the receiving agent's call balance in prior of making the call. These
balances are managed by a special public agent, named *Agent Zero* and are propagated through
the network. Withdraws of funds from the call balances is possible, however it is slower a operation, 
as it requires full propagation of the updated balance to prevent double spending.  

#### Invocations

Invocation commands represent the direct communication between agents. Invocation command encapsulate
a specific message exchanged between two agents and the respective object capability. For the underlying
decentralized network, the message is opaque object. Therefore it is up to the communicating agents to
agree upon a protocol and serialization mechanism. The messages are also the mechanism for propagating
state changes across the network and also carriers of object capabilities.

### Consensus

From consensus standpoint, Apocryph is proof-of-work network augmented with leaderless Byzantine fault
tolerance protocol inspired by [Wavelet](https://wavelet.perlin.net/whitepaper.pdf) and [Himitsu](https://www.youtube.com/watch?v=C542HhQKzKQ). 

Every group of agents is operated by a separate blockchain with a separate consensus that runs in 
parallel with other's blockchain consensus instances. Every consensus instance runs over a 
number of dedicated virtual nodes with different roles (proposers or validators). The process of 
selecting virtual nodes, the respective role and forming the proof-of-work network is called *selection*.

Another building block of Apocryph protocol are *facts*. Facts can be: slot claims, confirmed blocks,
rejected block or externally signed blocks. The high level goal of the protocol is to enable virtual nodes
to find the "ground truth" by combining their knowledge with observations of the facts produced
by other nodes. Virtual nodes gather knowledge by using *querying* other virtual nodes and validating their
responses. The observation of facts by a virtual node is achieved with *gossiping*.  

#### Selection

Virtual nodes selection is done using algorithm inspired by [Automaton](https://automaton.network/#i_koh)'s King of the Hill. 
For every consensus instance there is matrix of slots identified by bytes prefix. For every slot there is a random salt that changes 
periodically. Every node that participates in the network tries to produce private key the corresponds to 
a public key with a bytes prefix using an algorithm similar to [Vanitygen](https://en.bitcoin.it/wiki/Vanitygen). Upon discovery of a pair
of public and private keys with the desired properties the following procedure is initiated:

1. Prefix of public key is used for selecting a slot
2. The suffix of the public key is combined with the salt and hashed to produce the difficulty of the current pair.
3. Slot claim is produced using:
    - the public key;
    - the corresponding difficulty;
    - role (proposer or validator), randomly selected from the slots with the oldest salt. 
4. Gossip of the slot claim is distributed across the network.
5. The virtual node is selected based on the highest difficulty for the slot.

The selection algorithm has the following properties:

1. If different KotH networks share the same cryptography they can have shared miners
2. Keys with lower difficulty for a specific salt, might have higher difficulty for another salt,
so nodes have incentive to store the keys.  

#### Querying

Virtual nodes queries each other to reach consensus for the next block proposal. The querying
mechanism is following [Snowball](https://ipfs.io/ipfs/QmUy4jh5mGNZvLkjies1RWM4YuvJh5o2FYopNPVYwrRVGV) algorithm:

1. If a node is proposer, it generates block proposal and responds with it to all queries.
2. All virtual nodes start [Snowball](https://ipfs.io/ipfs/QmUy4jh5mGNZvLkjies1RWM4YuvJh5o2FYopNPVYwrRVGV) queries to each other and respond with:
    - their opinion (if proposers) OR
    - accept the query block proposal if it comes from a higher proposer (proposer preceeding current proposer in the proposers list) OR
    - the block proposal that they already have
3. In parallel with (2) all block proposals that are received from queries get validated (in priority queue)
4. Virtual node *commit* on the block proposal to which [Snowball](https://ipfs.io/ipfs/QmUy4jh5mGNZvLkjies1RWM4YuvJh5o2FYopNPVYwrRVGV) have converged.

In addition to the well known properties, Apocryph blocks have the following structure:

1. State
2. Commands
3. Object Capabilities
4. History (validators signatures from the last N valid blocks) 

#### Gossiping

Every virtual node, validates the committed block and based on the validation result the block is either *confirmed* or *rejected* (if invalid).
When a virtual node confirms or rejects a block it generates a fact and starts gossiping it. Every validator, gossips only facts
that are consistent with his observations. When 2/3 of virtual nodes sign a gossip it gets accepted. If the accepted gossip
is not consistent with the virtual node observations (ex. it is invalid), then the virtual node doesn't accept the block and forks. 

#### Agent Zero

Apocryph network contains one well known agent, called Agent Zero that generates the economic environment. Agent Zero is 
a regular agent that holds the following information in its state:

1. Main token balances
2. Agents call balances

It also responds to the following inter-agent communication messages:

1. Creating agents / chains
2. Deposit / withdraw of funds to call balances

Agent Zero chain (also referred as *main chain*) has the same consensus as any other chain. All agents include reference to the 
last block of the main chain in their blocks to maintain unified economic model.

#### Inter Blockchain Communication

Gossiping exchanges facts across the whole network. Based on the gossips, the virtual nodes observe not only their chain, 
but also other chains in the network. When a virtual node with proposer role, observes a fact containing accepted block
with invocation to his chain, the virtual node includes the invocation command in the next proposal for the current chain, using the following process:

1. It validates the provided capabilities in the block from the other chain.
2. It checks if call balance is sufficient, by referring to the last known block on the main chain. 
3. Generates new block proposal, using the state from the last accepted block on the current chain.

### Network Nodes

Apocryph is built on top of [Peprer](https://github.com/obecto/perper) - stream-based, horizontally 
scalable framework for asynchronous data processing. This enable Apocryph to run on physical nodes 
with various sizes: from single machine (using docker-compose) to a datacenter grade cluster environment
using [Kubernetes](http://kubernetes.io/). All physical (network) nodes form a decentralized network and communicate
with each other using [IPFS](https://ipfs.io/).

Every network node has the following components running on it:

1. Runtime
2. Agents
3. Client
4. Services

All of the components runs in separate containers, where the components running user code (Agents and Services) run in sandboxed 
container environments usign [gVisor](https://gvisor.dev/) 

#### Clients and Services

Client and Services are built-in mechanisms in Apocryph for interaction with the external world. Client exposes a Websocket
that allows thin clients (like mobile apps) to connect and observe Apocryph environment by watching for certain facts that are gossiped.
Clients enable also end-users to interact with the Apocryph environment by enabling them to gossip a specific facts (ex. deposit of funds) that might
be picked up by the proposers if they are valid. 

Services can be considered as automated clients, that directly watch and interact with the Apocryph environment by observing / gossiping specific facts, 
typically to assist agents with specific operations (ex. training a neural network).

#### Availability

Separation between network nodes and virtual nodes, enable high availability of the network. If network node, becomes unavailable, 
then another network node (that has mined the same key) can host the virtual node for the specific slot. 

In case of temporal network partition, the two networks will progress independently and lately converge as part of the querying process. 

### Test Harness

Using Docker Compose to run Apocryph runtime is the recommended way for users that
would like to run Apocryph Developer Node.

##### Prerequisite
- Install [Docker](https://docs.docker.com/install/)
- Install [Docker Compose](https://docs.docker.com/compose/install/)

##### Start IPFS Daemon

Apocryph uses IPFS for its DPoS consensus implementation, thus requires IPFS daemon to run locally on the node:

```bash
docker-compose up -d ipfs
```

##### Start Apocryph Runtime

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

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

#### Prerequisite

Before running this sample, you must have the following:

- The recommended operating system is Ubuntu 18.04 LTS.
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
