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
            var peerConnector = (new FakePeerConnectorProvider()).GetPeerConnector();

            var chain = new Chain(new ChainState(new MerkleTreeNode<AgentState>(new Hash<IMerkleTree<AgentState>>[] { }), 0), "", null, slotsCount);
            var chainId = await hashResolver.StoreAsync(chain);

            var cancellationTokenSource = new CancellationTokenSource();

            var kothStateStream = new FakeStream<(Hash<Chain>, Slot?[])>(KoTH.KoTHProcessor(null, new FakeState(), hashResolver, peerConnector, null, cancellationTokenSource.Token));

            var minerTask = SimpleMiner.Miner((kothStateStream, new Hash<Chain>[] { chainId }), hashResolver, peerConnector, cancellationTokenSource.Token);

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