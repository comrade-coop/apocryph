using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Channels;
using System.Threading.Tasks;
using Apocryph.Consensus;
using Apocryph.Ipfs;
using Apocryph.Ipfs.Fake;
using Apocryph.Ipfs.MerkleTree;
using Microsoft.Azure.WebJobs;
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
            var selfPeer = new Peer(Hash.From(0).Bytes);
            var hashResolver = new FakeHashResolver();

            var chain = new Chain(new ChainState(new MerkleTreeNode<AgentState>(new Hash<IMerkleTree<AgentState>>[] { }), 0), "", null, slotsCount);
            var chainId = await hashResolver.StoreAsync(chain);

            var tokenSource = new CancellationTokenSource();
            var minedKeysCollectorStream = new FakeAsyncCollector<(Hash<Chain>, Slot)>();
            var kothStateStream = KoTH.Processor(minedKeysCollectorStream, new FakeState(), hashResolver);

            kothStateStream = kothStateStream.Select(x => (x.Item1, x.Item2.ToArray())); // Duplicate the array, as KoTH modifies it by reference

            await minedKeysCollectorStream.AddAsync((chainId, new Slot(selfPeer, new byte[] { 0 })));

            var minerTask = SimpleMiner.Miner(("-", kothStateStream, selfPeer), minedKeysCollectorStream, tokenSource.Token);

            var i = 0;
            await foreach (var (stateChainId, peers) in kothStateStream)
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
            minedKeysCollectorStream.Complete();
            await minerTask;
        }

        public class FakeAsyncCollector<T> : IAsyncCollector<T>, IAsyncEnumerable<T> // FIXME: move to Perper
        {
            private Channel<T> _channel = Channel.CreateUnbounded<T>();

            public async Task AddAsync(T item, CancellationToken cancellationToken = default)
            {
                await _channel.Writer.WriteAsync(item, cancellationToken);
            }

            public Task FlushAsync(CancellationToken cancellationToken = default)
            {
                return Task.CompletedTask;
            }

            public void Complete(Exception? exception = null)
            {
                _channel.Writer.Complete(exception);
            }

            public IAsyncEnumerator<T> GetAsyncEnumerator(CancellationToken cancellationToken = default)
            {
                return _channel.Reader.ReadAllAsync().GetAsyncEnumerator(cancellationToken);
            }
        }
    }
}