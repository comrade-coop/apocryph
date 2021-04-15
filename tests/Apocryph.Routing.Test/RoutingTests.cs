using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Apocryph.Consensus;
using Apocryph.Ipfs;
using Apocryph.Ipfs.Fake;
using Apocryph.Ipfs.MerkleTree;
using Apocryph.Ipfs.Test;
using Apocryph.KoTH;
using Perper.WebJobs.Extensions.Fake;
using Perper.WebJobs.Extensions.Model;
using Xunit;

namespace Apocryph.Routing.Test
{
    public class RouterTests
    {
        private Message GenerateMessage(Hash<Chain> targetChain, int i, int nonce = 0)
        {
            return new Message(new Reference(targetChain, nonce, new[] { typeof(int).FullName! }), ReferenceData.From(i));
        }

        private async Task<Hash<Chain>> GetTestChain(IHashResolver hashResolver, string consensusType = "FakeConsensus")
        {
            var chain = new Chain(new ChainState(new MerkleTreeNode<AgentState>(new Hash<IMerkleTree<AgentState>>[] { }), 0), consensusType, null, 1);
            return await hashResolver.StoreAsync(chain);
        }

        private async Task RegisterConsensus(
            FakeAgent routingAgent, IHashResolver hashResolver,
            Hash<Chain> chainId,
            Func<
                (IAsyncEnumerable<Message> messages, string subscribtionsStream, Chain chain, IAsyncEnumerable<(Hash<Chain>, Slot?[])> kothStates, IAgent executor),
                Task<IAsyncEnumerable<Message>>> consensus)
        {
            var consensusType = (await hashResolver.RetrieveAsync(chainId)).ConsensusType;
            var consensusAgent = new FakeAgent();
            consensusAgent.RegisterFunction(consensusType, consensus);
            routingAgent.RegisterAgent(consensusType, () => consensusAgent);
        }

        private Task<(string, IStream<Message>)> InvokeConsensus(FakeAgent routingAgent, Hash<Chain> chainId)
        {
            return routingAgent.CallFunctionAsync<(string, IStream<Message>)>("GetChainInstance", chainId);
        }

        [Fact]
        public async void Routing_InstancesConsensus_PassingParams()
        {
            var hashResolver = new FakeHashResolver();
            var chainId = await GetTestChain(hashResolver);

            var executor = new FakeAgent();
            var kothStates = new (Hash<Chain>, Slot?[])[]
            {
                (chainId, new Slot?[] { new Slot(new Peer(new byte[] { 0 }), new byte[] { 0 }) }),
            };

            var routingAgent = await RoutingFakes.GetRoutingAgent(hashResolver, kothStates.ToAsyncEnumerable(), executor);

            await RegisterConsensus(routingAgent, hashResolver, chainId, async input =>
            {
                Assert.NotNull(input.messages);
                Assert.NotNull(input.subscribtionsStream);
                Assert.Equal(kothStates, await input.kothStates.ToArrayAsync(), SerializedComparer.Instance);
                Assert.Equal(chainId, Hash.From(input.chain));
                Assert.Equal(executor, input.executor);

                return AsyncEnumerable.Empty<Message>();
            });
            await InvokeConsensus(routingAgent, chainId);
        }

        [Fact]
        public async void Routing_Routes_Subscriptions()
        {
            var hashResolver = new FakeHashResolver();
            var routingAgent = await RoutingFakes.GetRoutingAgent(hashResolver);
            var chainIdFrom = await GetTestChain(hashResolver, "FromConsensus");
            var chainIdTo = await GetTestChain(hashResolver, "ToConsensus");

            var testMessages = Enumerable.Range(0, 10).Select(x => GenerateMessage(chainIdFrom, x, -1)).ToArray();
            var receivedMessages = new TaskCompletionSource<IAsyncEnumerable<Message>>();
            var subscribeReference = testMessages[0].Target;
            var fromInvoked = false;

            await RegisterConsensus(routingAgent, hashResolver, chainIdFrom, input =>
            {
                Assert.False(fromInvoked);
                fromInvoked = true;
                return Task.FromResult(testMessages.ToAsyncEnumerable());
            });

            await RegisterConsensus(routingAgent, hashResolver, chainIdTo, async input =>
            {
                Assert.False(fromInvoked);
                await FakeAgent.WriteToBlankStream(input.subscribtionsStream, new List<Reference> { subscribeReference });
                receivedMessages.SetResult(input.messages);
                return AsyncEnumerable.Empty<Message>();
            });

            await InvokeConsensus(routingAgent, chainIdTo);

            var enumerator = (await receivedMessages.Task).GetAsyncEnumerator();
            foreach (var message in testMessages)
            {
                Assert.True(await enumerator.MoveNextAsync());
                Assert.Equal(message, enumerator.Current, SerializedComparer.Instance);
            }
        }

        [Fact]
        public async void Routing_Routes_Calls()
        {
            var hashResolver = new FakeHashResolver();
            var routingAgent = await RoutingFakes.GetRoutingAgent(hashResolver);
            var chainIdFrom = await GetTestChain(hashResolver, "FromConsensus");
            var chainIdTo = await GetTestChain(hashResolver, "ToConsensus");

            var testMessages = Enumerable.Range(0, 10).Select(x => GenerateMessage(chainIdTo, x)).ToArray();
            var receivedMessages = new TaskCompletionSource<IAsyncEnumerable<Message>>();
            var toInvoked = false;

            await RegisterConsensus(routingAgent, hashResolver, chainIdFrom, input =>
            {
                Assert.False(toInvoked);
                return Task.FromResult(testMessages.ToAsyncEnumerable());
            });

            await RegisterConsensus(routingAgent, hashResolver, chainIdTo, input =>
            {
                Assert.False(toInvoked);
                toInvoked = true;
                receivedMessages.SetResult(input.messages);
                return Task.FromResult(AsyncEnumerable.Empty<Message>());
            });

            await InvokeConsensus(routingAgent, chainIdFrom);

            var enumerator = (await receivedMessages.Task).GetAsyncEnumerator();
            foreach (var message in testMessages)
            {
                Assert.True(await enumerator.MoveNextAsync());
                Assert.Equal(message, enumerator.Current, SerializedComparer.Instance);
            }
        }
    }
}