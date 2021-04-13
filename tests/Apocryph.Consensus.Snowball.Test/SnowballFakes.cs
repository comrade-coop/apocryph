using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Apocryph.Consensus.Snowball.FunctionApp;
using Apocryph.Ipfs;
using Apocryph.Ipfs.Fake;
using Apocryph.KoTH;
using Perper.WebJobs.Extensions.Fake;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Consensus.Snowball.Test
{
    public static class SnowballFakes
    {
        public static FakeAgent GetSnowballAgent(IHashResolver hashResolver, IPeerConnector peerConnector)
        {
            var agent = new FakeAgent();
            var snowball = new SnowballConsensus(new FakeContext(agent), new FakeState());

            agent.RegisterFunction("Apocryph-SnowballConsensus", ((IAsyncEnumerable<Message>, string, Chain, IAsyncEnumerable<(Hash<Chain>, Slot?[])>, IAgent) input) => snowball.Start(input, hashResolver));
            agent.RegisterFunction("SnowballStream", ((Hash<Chain>, Hash<Block>, SnowballParameters, IAgent) input) => snowball.SnowballStream(input, hashResolver, peerConnector).Select(x => x));
            agent.RegisterFunction("MessagePool", ((IAsyncEnumerable<Message>, Hash<Chain>) input) => snowball.MessagePool(input));
            agent.RegisterFunction("KothProcessor", ((Hash<Chain>, IAsyncEnumerable<(Hash<Chain>, Slot?[])>) input) => snowball.KothProcessor(input));

            return agent;
        }

        public static Task<IAsyncEnumerable<Message>> StartSnowballAgent((IAsyncEnumerable<Message> messages, string subscribtionsStream, Chain chain, IAsyncEnumerable<(Hash<Chain>, Slot?[])> kothStates, IAgent executor) input, IHashResolver hashResolver, IPeerConnector peerConnector)
        {
            var agent = GetSnowballAgent(hashResolver, peerConnector);

            return agent.CallFunctionAsync<IAsyncEnumerable<Message>>("Apocryph-SnowballConsensus", input);
        }

        public static Task<IAsyncEnumerable<Message>[]> StartSnowballNetwork(Peer[] peers, IHashResolver hashResolver, FakePeerConnectorProvider peerConnectorProvider, Chain chain, IAgent executor, Func<Peer, (IAsyncEnumerable<Message> messages, string subscribtionsStream, IAsyncEnumerable<(Hash<Chain>, Slot?[])> kothStates)> inputs)
        {
            var tasks = peers.Select(async peer =>
            {
                var peerConnector = peerConnectorProvider.GetConnector(peer);

                var (messages, subscribtionsStream, kothStates) = inputs(peer);

                return await StartSnowballAgent((messages, subscribtionsStream, chain, kothStates, executor), hashResolver, peerConnector);
            });

            return Task.WhenAll(tasks);
        }
    }
}