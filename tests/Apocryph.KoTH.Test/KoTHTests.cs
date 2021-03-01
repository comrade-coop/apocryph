using System.Linq;
using System.Threading.Channels;
using System.Threading.Tasks;
using Apocryph.Consensus;
using Apocryph.HashRegistry;
using Apocryph.HashRegistry.MerkleTree;
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

        [Fact]
        public async void KoTH_KeepsTrack_OfMinedPeers()
        {
            var hashRegistry = GetHashRegistryProxy();

            var chain = new Chain(new MerkleTreeNode<AgentState>(new Hash<IMerkleTree<AgentState>>[] { }), "Apocryph-DummyConsensus", 100);
            var chainId = await hashRegistry.StoreAsync(chain);

            var minedKeysChannel = Channel.CreateUnbounded<(Hash<Chain>, Peer)>();
            var outputStream = KoTH.Processor((minedKeysChannel.Reader.ReadAllAsync(), hashRegistry), new FakeState());

            var writerTask = Task.Run(async () =>
            {
                for (var i = 0; i < 200; i++)
                {
                    await minedKeysChannel.Writer.WriteAsync((chainId, new Peer(i, new byte[] { 0 })));
                }
                minedKeysChannel.Writer.Complete();
            });

            var previousCount = 0;
            await foreach (var (stateChainId, peers) in outputStream)
            {
                var count = peers.Count(x => x != null);
                Assert.Equal(stateChainId, chainId);
                Assert.True(count == previousCount || count == previousCount + 1);
                previousCount = count;
            }
            await writerTask;

        }
    }
}