using System.Collections.Generic;
using System.Linq;
using Apocryph.HashRegistry;
using Apocryph.PerperUtilities;
using Perper.WebJobs.Extensions.Fake;
using Xunit;

namespace Apocryph.Peering.Test
{
    public class PeeringTests
    {
        [Fact]
        public async void Connect_WithRegisteredPeer_EstablishesConnection()
        {
            var peerA = new Peer(Hash.From(0).Cast<object>());
            var peerB = new Peer(Hash.From(1).Cast<object>());

            var peerings = await PeeringFakes.GetPeeringAgents(new[] { peerA, peerB });
            var peeringA = peerings[0];
            var peeringB = peerings[1];

            var messagesIn = new object[] { 0, "123" };
            var messagesOut = new object[] { 1, "234" };
            var connectionType = "testConnection";

            var peerAgentA = new FakeAgent();
            peerAgentA.RegisterFunction("ConnectionHandler", async ((Peer sender, IAsyncEnumerable<object> messages) input) =>
            {
                Assert.Equal(input.sender, peerB);
                Assert.Equal(await input.messages.ToArrayAsync(), messagesIn);
                return messagesOut.ToAsyncEnumerable();
            });
            var handlerA = new Handler<IAsyncEnumerable<object>>(peerAgentA, "ConnectionHandler");

            await peeringA.CallActionAsync("Register", (connectionType, handlerA));

            var output = await peeringB.CallFunctionAsync<IAsyncEnumerable<object>>("Connect", (peerA, connectionType, messagesIn.ToAsyncEnumerable()));

            Assert.Equal(await output.ToArrayAsync(), messagesOut);
        }
    }
}