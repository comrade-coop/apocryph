using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Apocryph.Consensus;
using Apocryph.Ipfs;
using Apocryph.KoTH;
using Apocryph.Routing.FunctionApp;
using Perper.WebJobs.Extensions.Fake;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Routing.Test
{
    using Routing = Apocryph.Routing.FunctionApp.Routing;

    public static class RoutingFakes
    {
        public static async Task<FakeAgent> GetRoutingAgent(IHashResolver hashResolver, IAsyncEnumerable<(Hash<Chain>, Slot?[])>? kothStates = null, IAgent? executor = null)
        {
            var agent = new FakeAgent();
            var state = new FakeState();
            var context = new FakeContext(agent);

            agent.RegisterFunction("GetChainInstance", (Hash<Chain> input) => Routing.GetChainInstance(input, context, state, hashResolver));
            agent.RegisterFunction("RouterInput", async ((IAsyncEnumerable<Message>, IAsyncEnumerable<List<Reference>>) input) =>
                RouterInput.RunAsync(input, context, await ((IState)state).Entry<List<Reference>?>(Guid.NewGuid().ToString(), () => null)).Select(x => x));
            agent.RegisterFunction("RouterOutput", ((string publicationsStream, IAsyncEnumerable<Message>, Hash<Chain>) input) =>
                RouterOutput.RunAsync(input, FakeAgent.GetBlankStreamCollector<Message>(input.publicationsStream), context, state));
            agent.RegisterFunction("PostMessage", ((string stream, Message message) input) => FakeAgent.WriteToBlankStream(input.stream, input.message));

            await Routing.Start((kothStates ?? AsyncEnumerable.Empty<(Hash<Chain>, Slot?[])>(), executor ?? new FakeAgent()), state);

            return agent;
        }
    }
}