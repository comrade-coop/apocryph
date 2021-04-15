using System;
using System.Collections.Generic;
using System.Threading.Tasks;
using Apocryph.Consensus;
using Apocryph.Consensus.Snowball;
using Apocryph.Ipfs;
using Apocryph.Ipfs.MerkleTree;
using Apocryph.KoTH;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Model;
using Perper.WebJobs.Extensions.Triggers;
using SampleAgents.FunctionApp.Agents;

namespace SampleAgents.FunctionApp
{
    public static class Launcher
    {
        [FunctionName("SampleAgents")]
        public static async Task Start([PerperTrigger] object? input, IContext context, IHashResolver hashResolver, IPeerConnector peerConnector)
        {
            var (executorAgent, _) = await context.StartAgentAsync<object?>("Apocryph-Executor", null);

            await executorAgent.CallActionAsync("Register", (Hash.From("AgentOne"), context.Agent, "AgentOne"));
            await executorAgent.CallActionAsync("Register", (Hash.From("AgentTwo"), context.Agent, "AgentTwo"));

            var agentStates = new[] {
                new AgentState(0, ReferenceData.From(new AgentOne.AgentOneState()), Hash.From("AgentOne")),
                new AgentState(1, ReferenceData.From(new AgentTwo.AgentTwoState()), Hash.From("AgentTwo"))
            };

            var agentStatesTree = await MerkleTreeBuilder.CreateRootFromValues(hashResolver, agentStates, 2);

            Chain chain;

            if (Environment.GetEnvironmentVariable("SAMPLE_AGENTS_CONSENSUS") == "Dummy")
            {
                chain = new Chain(new ChainState(agentStatesTree, agentStates.Length), "Apocryph-DummyConsensus", null, 1);
            }
            else
            {
                var snowballParameters = await hashResolver.StoreAsync<object>(new SnowballParameters(15, 0.8, 25));
                chain = new Chain(new ChainState(agentStatesTree, agentStates.Length), "Apocryph-SnowballConsensus", snowballParameters, 60);
            }

            var chainId = await hashResolver.StoreAsync(chain);

            var (_, kothStates) = await context.StartAgentAsync<IAsyncEnumerable<(Hash<Chain>, Slot?[])>>("Apocryph-KoTH", null);
            var minerAgentTask = context.StartAgentAsync<object?>("Apocryph-KoTH-SimpleMiner", (kothStates, new Hash<Chain>[] { chainId }));

            var (routingAgent, _) = await context.StartAgentAsync<object?>("Apocryph-Routing", (kothStates, executorAgent));

            var (chainInput, chainOutputs) = await routingAgent.CallFunctionAsync<(string, IStream<Message>)>("GetChainInstance", chainId);

            System.Console.WriteLine(chainInput);

            await routingAgent.CallActionAsync("PostMessage", (chainInput, new Message(
                new Reference(chainId, 0, new[] { typeof(PingPongMessage).FullName! }),
                ReferenceData.From(new PingPongMessage(
                    callback: new Reference(chainId, 1, new[] { typeof(PingPongMessage).FullName! }),
                    content: "START! ",
                    accumulatedValue: 0
                ))
            )));
        }
    }
}