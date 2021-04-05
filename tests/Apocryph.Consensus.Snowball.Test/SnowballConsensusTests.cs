using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Executor.Test;
using Apocryph.HashRegistry;
using Apocryph.HashRegistry.Test;
using Apocryph.KoTH;
using Apocryph.Peering.Test;
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
            var executor = await ExecutorFakes.GetExecutor(ExecutorFakes.TestAgents);
            var peers = PeeringFakes.GetFakePeers(peersCount);
            var peeringRouter = PeeringFakes.GetPeeringRouter();

            var snowballParameters = new SnowballParameters(3, 0.5, 25);
            var (chain, inputMessages, expectedOutputMessages) = await ExecutorFakes.GetTestAgentScenario(hashRegistry, "Apocryph-SnowballConsensus", snowballParameters, peersCount);

            var chainId = Hash.From(chain);

            var kothStates = new (Hash<Chain>, Slot?[])[]
            {
                (chainId, peers.Select(peer => (Slot?)new Slot(peer, new byte[0])).ToArray())
            };

            var outputStreams = await SnowballFakes.StartSnowballNetwork(peers, hashRegistry, chain, peeringRouter, executor,
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