using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Executor.Test;
using Apocryph.Ipfs;
using Apocryph.Ipfs.Fake;
using Apocryph.Ipfs.Test;
using Apocryph.KoTH;
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
            var hashResolver = new FakeHashResolver();
            var executor = await ExecutorFakes.GetExecutor(ExecutorFakes.TestAgents);
            var peerConnectorProvider = new FakePeerConnectorProvider();
            var peers = Enumerable.Range(0, peersCount).Select(x => FakePeerConnectorProvider.GetFakePeer()).ToArray();

            var snowballParameters = new SnowballParameters(3, 0.5, 25);
            var (chain, inputMessages, expectedOutputMessages) = await ExecutorFakes.GetTestAgentScenario(hashResolver, "Apocryph-SnowballConsensus", snowballParameters, peersCount);

            var chainId = Hash.From(chain);

            var started = new TaskCompletionSource<bool>();
            async IAsyncEnumerable<(Hash<Chain>, Slot?[])> kothStates()
            {
                await started.Task;
                yield return (chainId, peers.Select(peer => (Slot?)new Slot(peer, new byte[0])).ToArray());
            }

            var outputStreams = await SnowballFakes.StartSnowballNetwork(peers, hashResolver, peerConnectorProvider, chain, executor,
                peer => (inputMessages.ToAsyncEnumerable(), "-", kothStates()));

            started.TrySetResult(true);

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