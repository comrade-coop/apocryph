using System;
using System.Linq;
using Apocryph.Consensus;
using Apocryph.Ipfs;
using Apocryph.Ipfs.Fake;
using Apocryph.Ipfs.MerkleTree;
using Perper.WebJobs.Extensions.Dataflow;
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
        public async void KoTH_KeepsTrack_OfMinedPeers(int slots, int mineCount)
        {
            var selfPeer = new Peer(Hash.From(0).Bytes);
            var hashResolver = new FakeHashResolver();

            var chain = new Chain(new ChainState(new MerkleTreeNode<AgentState>(new Hash<IMerkleTree<AgentState>>[] { }), 0), "", null, slots);
            var chainId = await hashResolver.StoreAsync(chain);

            var minedKeys = Enumerable.Range(0, mineCount).Select(i => (chainId, new Slot(selfPeer, BitConverter.GetBytes(i))));
            var outputStream = KoTH.Processor(minedKeys.ToAsyncEnumerable(), new FakeState(), hashResolver);

            var previousCount = 0;
            await foreach (var (stateChainId, peers) in outputStream)
            {
                var count = peers.Count(x => x != null);
                Assert.Equal(stateChainId, chainId);
                Assert.True(count == previousCount || count == previousCount + 1);
                previousCount = count;
            }
        }
    }
}