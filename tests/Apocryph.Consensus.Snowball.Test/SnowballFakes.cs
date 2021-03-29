using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Apocryph.Consensus.Snowball.FunctionApp;
using Apocryph.HashRegistry;
using Apocryph.KoTH;
using Apocryph.Peering;
using Perper.WebJobs.Extensions.Fake;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Consensus.Snowball.Test
{
    public static class SnowballFakes
    {
        public static FakeAgent GetSnowballAgent()
        {
            var agent = new FakeAgent();
            var snowball = new SnowballConsensus(new FakeContext(agent), new FakeState());

            agent.RegisterFunction("Apocryph-SnowballConsensus", ((IAsyncEnumerable<Message>, string, HashRegistryProxy, Chain, IAgent, IAsyncEnumerable<(Hash<Chain>, Slot?[])>) input) => snowball.Start(input));
            agent.RegisterFunction("SnowballStream", ((IAgent, HashRegistryProxy, SnowballParameters, Hash<Chain>, Hash<Block>) input) => snowball.SnowballStream(input).Select(x => x));
            agent.RegisterFunction("MessagePool", ((IAsyncEnumerable<Message>, Hash<Chain>) input) => snowball.MessagePool(input));
            agent.RegisterFunction("KothProcessor", ((Hash<Chain>, IAsyncEnumerable<(Hash<Chain>, Slot?[])>) input) => snowball.KothProcessor(input));
            agent.RegisterFunction("PeeringResponder", ((Peer, IAsyncEnumerable<object>) input) => snowball.PeeringResponder(input));

            return agent;
        }

        public static Task<IAsyncEnumerable<Message>> StartSnowballAgent(Action<FakeAgent> initSnowballAgent, (IAsyncEnumerable<Message> messages, string subscribtionsStream, HashRegistryProxy hashRegistry, Chain chain, IAgent peering, IAsyncEnumerable<(Hash<Chain>, Slot?[])> kothStates) input)
        {
            var agent = GetSnowballAgent();
            initSnowballAgent(agent);

            return agent.CallFunctionAsync<IAsyncEnumerable<Message>>("Apocryph-SnowballConsensus", input);
        }

        public static Task<IAsyncEnumerable<Message>[]> StartSnowballNetwork(Peer[] peers, HashRegistryProxy hashRegistry, Chain chain, IAgent peeringRouter, Action<FakeAgent> initSnowballAgent, Func<Peer, (IAsyncEnumerable<Message> messages, string subscribtionsStream, IAsyncEnumerable<(Hash<Chain>, Slot?[])> kothStates)> inputs)
        {
            var tasks = peers.Select(async peer =>
            {
                var peering = await peeringRouter.CallFunctionAsync<IAgent>("_GetPeering", peer);

                var (messages, subscribtionsStream, kothStates) = inputs(peer);

                return await StartSnowballAgent(initSnowballAgent, (messages, subscribtionsStream, hashRegistry, chain, peering, kothStates));
            });

            return Task.WhenAll(tasks);
        }
    }
}