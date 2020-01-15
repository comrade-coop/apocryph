[Getting Started](https://github.com/comrade-coop/apocryph/blob/master/GettingStarted.md) | [Discord Community](https://discord.gg/ESr9KMR) 

----------------------------------

# APOCRYPH
## Consensus Network for Autonomous Agents

> Apocryph agents can automate the cash flow in autonomous organizations, optimize city traffic, or reward the computing power used to train neural networks. 

As engineers, we strive to automate everything. For us, the biggest promise of blockchain technology is that it can for the first time enable fully automatic and thus incorruptible social institutions. We can establish programmatic organizations and even whole programmatic economies that have the potential to drive our civilization to a new and unprecedented level of collaboration and growth. 

To build these new economies, developers need mature languages and scalable runtimes, which are still not available in mainstream blockchain networks. This motivated us to take a different approach and design a blockchain network that reuses as many established technologies as possible, instead of rewriting everything from scratch. As a result, we have built Apocryph - a consensus network for autonomous agents with the following advantages:

## Developer productivity
Apocryph is built on top of [Perper](https://github.com/obecto/perper) - a serverless stream processing framework maintained separately by members of our team. The main entities in Apocryph are called Agents - they process incoming messages as a stream and send new ones as output. Each Agent runs as a containerized lambda function app, that can be written in C#, Python, JavaScript or come as a WebAssembly.

## Proactive entities
Since autonomy is our main focus, Agents can proactively initiate their own execution by scheduling reminders and subscribing to events from other Agents. This opens entirely new use cases and more natural programming models. For example, users can have their own Proxy Agents on the network to actively manage and monitor tasks for them.

## Free user transactions
Messages in the network are processed only if they come with a valid execution ticket for the computing resources requested. These tickets can be paid either by the sending or the receiving party, while agents have their own wallets so they can pay for executing messages coming from certain users. This is crucial to enable use cases like voting, rating and user feedback that were unfeasible with previous blockchain economic models that require users to pay each transaction.

## Scalable network
The state of each Agent is stored on a separate blockchain and the transactions on these blockchains are running in parallel to enable an extremely high transaction throughput. Agents are self-governing and each agent can declare which subset of the Apocryph network validators should validate the respective Agent’s blockchain. 

## Scalable nodes
Since the network nodes run on Perper, they scale horizontally and each node is supposed to be more a cluster of machines, rather than just a single machine. This enables the network validators to run nodes on professional infrastructure and achieve economies of scale. 

## Interoperability and extensibility
With the main focus on decoupling and reusability, the Apocryph network is designed to be highly interoperable and extensible. Agents might require that certain services are running on their validator nodes and thus access any functionality needed - from sending emails to interoperability with EthereumEthereum interoperability and decentralized training of AI.

## DPoS consensus
The network runs on a Delegated Proof of Stake, Practical Byzantine Fault Tolerance consensus algorithm. IPFS is used for storing the block data, while block propagation and voting happen over IPFS GossipSub. By reusing established technology and existing infrastructure we significantly simplify the consensus engine, allowing more idiomatic and secure implementation.

## Built by a coop
Apocryph is built by the [Comrade Cooperative](https://www.comrade.coop/) - a member-owned organization of software developers and innovation builders, that is based on transparency, technocracy, and self-governance. In the past two years, we are working on two pillar projects around the most important use cases we saw for consensus networks - autonomous organizations with [Wetonomy](https://www.wetonomy.com/) and decentralized AI with [ScyNet](https://www.scynet.ai/). Apocryph emerged as a solution to the numerous problems we encountered while we were working on these projects and now they all form a coherent ecosystem.
