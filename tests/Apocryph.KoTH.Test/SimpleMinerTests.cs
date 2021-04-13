using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Consensus;
using Apocryph.Ipfs;
using Apocryph.Ipfs.Fake;
using Apocryph.Ipfs.MerkleTree;
using Perper.WebJobs.Extensions.Fake;
using Xunit;

namespace Apocryph.KoTH.Test
{
    using KoTH = Apocryph.KoTH.FunctionApp.KoTH;
    using SimpleMiner = Apocryph.KoTH.SimpleMiner.FunctionApp.SimpleMiner;

    public class SimpleMinerTests
    {
        [Theory]
        [InlineData(1)]
        [InlineData(10)]
#if SLOWTESTS
        [InlineData(100)]
#endif
        public async void SimpleMiner_Fills_AllPeers(int slotsCount)
        {
            var hashResolver = new FakeHashResolver();
            var peerConnector = (new FakePeerConnectorProvider()).GetConnector();

            var chain = new Chain(new ChainState(new MerkleTreeNode<AgentState>(new Hash<IMerkleTree<AgentState>>[] { }), 0), "", null, slotsCount);
            var chainId = await hashResolver.StoreAsync(chain);

            var cancellationTokenSource = new CancellationTokenSource();

            var kothStateStream = await KoTH.Processor(null, new FakeState(), hashResolver, peerConnector, cancellationTokenSource.Token);

            var ready = new TaskCompletionSource<bool>(); // NOTE: Used since we want to start things only after the miner is listening
            async IAsyncEnumerable<(Hash<Chain>, Slot?[])> kothStates()
            {
                ready.SetResult(true);
                await foreach (var state in kothStateStream)
                    yield return state;
            }

            var minerTask = SimpleMiner.Miner(kothStates(), peerConnector, cancellationTokenSource.Token);

            await ready.Task;
            await peerConnector.SendPubSub(KoTH.PubSubPath, (chainId, new Slot(peerConnector.Self, new byte[] { 0 })));

            var i = 0;
            await foreach (var (stateChainId, peers) in kothStateStream)
            {
                i++;
                Assert.True(i < slotsCount * 10); // Hang prevention
                var count = peers.Count(x => x != null);
                if (count == slotsCount)
                {
                    break;
                }
            }
            cancellationTokenSource.Cancel();

            await Task.WhenAll(peerConnector.PendingHandlerTasks);
            await minerTask;
        }
    }
}