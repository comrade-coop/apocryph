using System;
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

    public class KoTHTests
    {
        [Theory]
        [InlineData(10, 20)]
#if SLOWTESTS
        [InlineData(100, 200)]
#endif
        public async void KoTH_KeepsTrack_OfMinedPeers(int slotsCount, int mineCount)
        {
            var hashResolver = new FakeHashResolver();
            var peerConnector = (new FakePeerConnectorProvider()).GetPeerConnector();

            var chain = new Chain(new ChainState(new MerkleTreeNode<AgentState>(new Hash<IMerkleTree<AgentState>>[] { }), 0), "", null, slotsCount);
            var chainId = await hashResolver.StoreAsync(chain);

            var cancellationTokenSource = new CancellationTokenSource();

            var outputStream = await KoTH.KoTHProcessor(null, new FakeState(), hashResolver, peerConnector, null, cancellationTokenSource.Token);

            var _ = Task.Run(async () =>
            {
                for (var i = 0; i < mineCount; i++)
                {
                    var slot = new Slot(peerConnector.Self, BitConverter.GetBytes(i));
                    await peerConnector.SendPubSub(KoTH.PubSubPath, (chainId, slot));
                }

                await Task.WhenAll(peerConnector.PendingHandlerTasks);

                cancellationTokenSource.Cancel();
            });

            var previousCount = 0;
            await foreach (var (stateChainId, slots) in outputStream)
            {
                var count = slots.Count(x => x != null);
                Assert.Equal(stateChainId, chainId);
                Assert.True(count == previousCount || count == previousCount + 1);
                previousCount = count;
            }
        }
    }
}