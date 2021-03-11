using System;
using System.Linq;
using Apocryph.Consensus;
using Apocryph.HashRegistry;
using Apocryph.HashRegistry.MerkleTree;
using Apocryph.HashRegistry.Test;
using Apocryph.Peering;
using Perper.WebJobs.Extensions.Dataflow;
using Perper.WebJobs.Extensions.Fake;
using Xunit;

namespace Apocryph.KoTH.Test
{
    using HashRegistry = Apocryph.HashRegistry.FunctionApp.HashRegistry;
    using KoTH = Apocryph.KoTH.FunctionApp.KoTH;

    public class KoTHTests
    {
        private HashRegistryProxy GetHashRegistryProxy()
        {
            var registry = new HashRegistry(new FakeState());

            var agent = new FakeAgent();
            agent.RegisterFunction("Store", (byte[] data) => registry.Store(data, default));
            agent.RegisterFunction("Retrieve", (Hash hash) => registry.Retrieve(hash, default));

            return new HashRegistryProxy(agent);
        }

        [Theory]
        [InlineData(10, 20)]
        //[InlineData(100, 200)]
        public async void KoTH_KeepsTrack_OfMinedPeers(int slots, int mineCount)
        {
            var selfPeer = new Peer(Hash.From(0).Cast<object>());
            var hashRegistry = HashRegistryFakes.GetHashRegistryProxy();

            var chain = new Chain(new MerkleTreeNode<AgentState>(new Hash<IMerkleTree<AgentState>>[] { }), "Apocryph-DummyConsensus", slots);
            var chainId = await hashRegistry.StoreAsync(chain);

            var minedKeys = Enumerable.Range(0, mineCount).Select(i => (chainId, new Slot(selfPeer, BitConverter.GetBytes(i))));
            var outputStream = KoTH.Processor((minedKeys.ToAsyncEnumerable(), hashRegistry), new FakeState());

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