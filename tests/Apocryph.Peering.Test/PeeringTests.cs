using System.Collections.Generic;
using System.Linq;
using Apocryph.HashRegistry;
using Perper.WebJobs.Extensions.Fake;
using Xunit;

namespace Apocryph.Peering.Test
{
    using Peering = Apocryph.Peering.FunctionApp.Peering;

    public class ServiceRegistryTests
    {
        [Fact]
        public async void Connect_WithRegisteredPeer_EstablishesConnection()
        {
            var peering = new Peering(new FakeState());

            var messagesIn = new object[] { 0, "123" };
            var messagesOut = new object[] { 1, "234" };
            var connectionType = "testConnection";

            var peer = new Peer(Hash.From(0).Cast<object>());

            var peerAgent = new FakeAgent();
            peerAgent.RegisterFunction("ConnectionHandler", async ((string connectionType, IAsyncEnumerable<object> messages) input) =>
            {
                Assert.Equal(input.connectionType, connectionType);
                Assert.Equal(await input.messages.ToArrayAsync(), messagesIn);
                return messagesOut.ToAsyncEnumerable();
            });
            var handler = new PeerHandler(peerAgent, "ConnectionHandler");

            await peering.Register((peer, handler));

            var output = await peering.Connect((peer, connectionType, messagesIn.ToAsyncEnumerable()));

            Assert.Equal(await output.ToArrayAsync(), messagesOut);
        }
    }
}