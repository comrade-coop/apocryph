using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Apocryph.HashRegistry;
using Apocryph.PerperUtilities;
using Perper.WebJobs.Extensions.Fake;

namespace Apocryph.Peering.Test
{
    using PeerHandler = Handler<IAsyncEnumerable<object>>;
    using Peering = Apocryph.Peering.FunctionApp.Peering;

    public class PeeringFakes
    {
        public static FakeAgent GetPeeringRouter()
        {
            var peering = new Peering(new FakeState());

            var agent = new FakeAgent();
            agent.RegisterFunction("Apocryph-Peering", (Peering.PeeringState? input) => peering.Start(input));
            agent.RegisterFunction("_GetPeering", (Peer input) => peering._GetPeering(input, new FakeContext(agent)));
            agent.RegisterFunction("_GetHandler", ((Peer, string) input) => peering._GetHandler(input));
            agent.RegisterFunction("_SetHandler", ((Peer, string, PeerHandler) input) => peering._SetHandler(input));
            agent.RegisterFunction("Connect", ((Peer, string, IAsyncEnumerable<object>) input) => peering.Connect(input));
            agent.RegisterFunction("Register", ((string, PeerHandler) input) => peering.Register(input));

            agent.RegisterAgent("Apocryph-Peering", GetPeeringRouter);

            return agent;
        }

        public static Task<FakeAgent[]> GetPeeringAgents(IEnumerable<Peer> peers)
        {
            var agent = GetPeeringRouter();

            return Task.WhenAll(peers.Select(x => agent.CallFunctionAsync<FakeAgent>("_GetPeering", x)));
        }

        public static Peer[] GetFakePeers(int count)
        {
            return Enumerable.Range(0, count).Select(x => new Peer(Hash.From(x).Cast<object>())).ToArray();
        }
    }
}