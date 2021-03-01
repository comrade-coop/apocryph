using System.Linq;
using System.Threading;
using System.Threading.Channels;
using System.Threading.Tasks;
using Apocryph.Consensus;
using Apocryph.HashRegistry;
using Apocryph.HashRegistry.MerkleTree;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Fake;
using Xunit;

namespace Apocryph.KoTH.Test
{
    using HashRegistry = Apocryph.HashRegistry.FunctionApp.HashRegistry;
    using KoTH = Apocryph.KoTH.FunctionApp.KoTH;
    using SimpleMiner = Apocryph.KoTH.SimpleMiner.FunctionApp.SimpleMiner;

    public class SimpleMinerTests
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
        [InlineData(1)]
        [InlineData(10)]
        [InlineData(100)]
        public async void SimpleMiner_Fills_AllPeers(int slotsCount)
        {
            var hashRegistry = GetHashRegistryProxy();

            var chain = new Chain(new MerkleTreeNode<AgentState>(new Hash<IMerkleTree<AgentState>>[] { }), "Apocryph-DummyConsensus", slotsCount);
            var chainId = await hashRegistry.StoreAsync(chain);

            var tokenSource = new CancellationTokenSource();
            var minedKeysChannel = Channel.CreateUnbounded<(Hash<Chain>, Peer)>();
            var outputStream = KoTH.Processor((minedKeysChannel.Reader.ReadAllAsync(), hashRegistry), new FakeState());

            outputStream = outputStream.Select(x => (x.Item1, x.Item2.ToArray())); // Duplicate the array, as it is not immutable otherwise..

            await minedKeysChannel.Writer.WriteAsync((chainId, new Peer(0, new byte[] { 0 })));

            var minerTask = SimpleMiner.Miner(("-", outputStream, 1), new ChannelAsyncCollector<(Hash<Chain>, Peer)>(minedKeysChannel.Writer), tokenSource.Token);

            var i = 0;
            await foreach (var (stateChainId, peers) in outputStream)
            {
                i++;
                Assert.True(i < slotsCount * 10); // Prevent hangs
                var count = peers.Count(x => x != null);
                if (count == slotsCount)
                {
                    break;
                }
            }
            tokenSource.Cancel();
            minedKeysChannel.Writer.Complete();
            await minerTask;
        }

        public class ChannelAsyncCollector<T> : IAsyncCollector<T>
        {
            public ChannelWriter<T> ChannelWriter { get; set; }

            public ChannelAsyncCollector(ChannelWriter<T> channelWriter)
            {
                ChannelWriter = channelWriter;
            }

            public async Task AddAsync(T item, CancellationToken cancellationToken = default)
            {
                await ChannelWriter.WriteAsync(item);
            }

            public Task FlushAsync(CancellationToken cancellationToken = default)
            {
                return Task.CompletedTask;
            }
        }
    }
}