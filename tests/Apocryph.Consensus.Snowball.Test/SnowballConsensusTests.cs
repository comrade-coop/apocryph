using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Consensus.Test;
using Apocryph.HashRegistry;
using Apocryph.HashRegistry.MerkleTree;
using Apocryph.HashRegistry.Test;
using Apocryph.KoTH;
using Apocryph.Peering.Test;
using Perper.WebJobs.Extensions.Fake;
using Xunit;
using Xunit.Abstractions;

namespace Apocryph.Consensus.Snowball.Test
{
    public class SnowballConsensusTests
    {
        private readonly ITestOutputHelper _output;
        public SnowballConsensusTests(ITestOutputHelper output)
        {
            _output = output;
        }

        [Theory]
        [InlineData(5)]
#if SLOWTESTS
        [InlineData(50)] // Test mocks have race conditions after ~60
#endif
        public async void Snowball_ConfirmsAndExecutes_SingleMessage(int peersCount)
        {
            var hashRegistry = HashRegistryFakes.GetHashRegistryProxy();
            var peers = PeeringFakes.GetFakePeers(peersCount);
            var peeringRouter = PeeringFakes.GetPeeringRouter();

            var agentStates = new[] {
                new AgentState(0, ReferenceData.From("123"), "Agent1")
            };

            var agentStatesTree = await MerkleTreeBuilder.CreateRootFromValues(hashRegistry, agentStates, 2);

            var snowballParameters = new SnowballParameters(3, 0.5, 25);

            var chain = new Chain(new ChainState(agentStatesTree, 1), "Apocryph-SnowballConsensus", snowballParameters, peersCount);
            var chainId = await hashRegistry.StoreAsync(chain);

            var testMessageAllowed = new string[] { typeof(int).FullName! };

            var inputMessages = new Message[]
            {
                new Message(new Reference(chainId, 0, testMessageAllowed), ReferenceData.From(4)),
                new Message(new Reference(chainId, 0, testMessageAllowed), ReferenceData.From(5))
            };

            var expectedOutputMessages = new Message[]
            {
                new Message(new Reference(chainId, 1, testMessageAllowed), ReferenceData.From(3)),
                new Message(new Reference(chainId, 1, testMessageAllowed), ReferenceData.From(4))
            };

            var kothStates = new (Hash<Chain>, Slot?[])[]
            {
                (chainId, peers.Select(peer => (Slot?)new Slot(peer, new byte[0])).ToArray())
            };

            var agent1 = new FakeAgent();
            agent1.RegisterFunction("Agent1", ((AgentState state, Message message) input) =>
            {
                Assert.Equal(input.state, agentStates[0], SerializedComparer.Instance);
                var result = input.message.Data.Deserialize<int>() - 1;
                return (input.state, new[] { new Message(new Reference(chainId, 1, testMessageAllowed), ReferenceData.From(result)) });
            });

            var outputStreams = await SnowballFakes.StartSnowballNetwork(peers, hashRegistry, chain, peeringRouter,
                agent =>
                {
                    agent.RegisterAgent("Agent1", () => agent1);
                },
                peer => (inputMessages.ToAsyncEnumerable(), "-", kothStates.ToAsyncEnumerable()));


            // NOTE: Using WhenAll here, since (A) we need to iterate all streams fully and (B) PerperFakeStream does not propagate completion

            var cancellation = new CancellationTokenSource();

            await Task.WhenAll(outputStreams.Select(async stream =>
            {
                try
                {
                    var enumerator = stream.GetAsyncEnumerator();
                    foreach (var expectedMessage in expectedOutputMessages)
                    {
                        await enumerator.MoveNextAsync(cancellation.Token);
                        Assert.Equal(enumerator.Current, expectedMessage, SerializedComparer.Instance);
                    }
                    // NOTE: Would be nice to confirm there are no other stream items
                }
                catch (Exception e)
                {
                    _output.WriteLine(e.ToString());
                    throw;
                }
            }));

            cancellation.Cancel();
        }
    }
}