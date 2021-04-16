using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Xunit;

namespace Apocryph.Ipfs.Test
{
    [Collection("Ipfs Collection")]
    public class PeerConnectorTests
    {
        private IpfsFixture _fixture;

        public PeerConnectorTests(IpfsFixture fixture)
        {
            _fixture = fixture;
        }

        public static IEnumerable<object?[]> SampleQueryData() =>
            IpfsFixture.PeerConnectorImplementations.SelectMany((i) =>
                TestData.DataInterface.Select((d) => new object?[] { i, d[0], d[0] }));

        public static IEnumerable<object?[]> SamplePubSubData() =>
            IpfsFixture.PeerConnectorImplementations.SelectMany((i) =>
                TestData.DataInterface.Select((d) => new object?[] { i, d[0] }));

        [Theory]
        [MemberData(nameof(SampleQueryData))]
        public async void Query_AfterListenQuery_SendsQueryAndResult(string connectorImplementation, IExample dataRequest, IExample dataReply)
        {
            var connectorTo = _fixture.GetPeerConnector(connectorImplementation, 1);
            var connectorFrom = _fixture.GetPeerConnector(connectorImplementation, 2);

            var cancellationTokenSource = new CancellationTokenSource();
            var path = $"test-{Guid.NewGuid()}";

            await connectorTo.ListenQuery<IExample, IExample>(path, async (otherPeer, request) =>
            {
                Assert.Equal(otherPeer, await connectorFrom.Self);
                Assert.Equal(request, dataRequest, SerializedComparer.Instance);

                return dataRequest;
            }, cancellationTokenSource.Token);

            var reply = await connectorFrom.SendQuery<IExample, IExample>(await connectorTo.Self, path, dataRequest, cancellationTokenSource.Token);
            Assert.Equal(reply, dataReply, SerializedComparer.Instance);

            cancellationTokenSource.Cancel();
        }

        [Theory]
        [MemberData(nameof(SampleQueryData))]
        public async void Query_ToSelf_IsStillReceived(string connectorImplementation, IExample dataRequest, IExample dataReply)
        {
            var connector = _fixture.GetPeerConnector(connectorImplementation, 1);

            var cancellationTokenSource = new CancellationTokenSource();
            var path = $"test-{Guid.NewGuid()}";

            await connector.ListenQuery<IExample, IExample>(path, async (otherPeer, request) =>
            {
                Assert.Equal(otherPeer, await connector.Self);
                Assert.Equal(request, dataRequest, SerializedComparer.Instance);

                return dataRequest;
            }, cancellationTokenSource.Token);

            var reply = await connector.SendQuery<IExample, IExample>(await connector.Self, path, dataRequest, cancellationTokenSource.Token);
            Assert.Equal(reply, dataReply, SerializedComparer.Instance);
            cancellationTokenSource.Cancel();
        }


#if SLOWTESTS
        [Theory] // (Skip="Flaky on IPFS")
#else
        [Theory]
#endif
        [MemberData(nameof(SamplePubSubData))]
        public async void PubSub_AfterListenPubSub_SendsPubSub(string connectorImplementation, IExample dataMessage)
        {
            var connectorTo = _fixture.GetPeerConnector(connectorImplementation, 1);
            var connectorFrom = _fixture.GetPeerConnector(connectorImplementation, 2);

            var cancellationTokenSource = new CancellationTokenSource(TestData.WaitTime * 3); // Prevent hangs

            var path = $"test-{Guid.NewGuid()}";

            var receiveTasks = new List<Task>();
            foreach (var connector in new[] { connectorTo, connectorFrom })
            {
                var received = new TaskCompletionSource<bool>();
                receiveTasks.Add(received.Task);
                cancellationTokenSource.Token.Register(() => received.TrySetCanceled());

                await connector.ListenPubSub<IExample>(path, async (otherPeer, message) =>
                {
                    if (message == null) return false;
                    Assert.Equal(otherPeer, await connectorFrom.Self);
                    Assert.Equal(message, dataMessage, SerializedComparer.Instance);
                    received.SetResult(true);
                    return true;
                }, cancellationTokenSource.Token);
                await connectorFrom.SendPubSub<IExample?>(path, null);
            };

            await Task.Delay(TestData.WaitTime);

            await connectorFrom.SendPubSub<IExample>(path, dataMessage, cancellationTokenSource.Token);

            await Task.WhenAll(receiveTasks);

            cancellationTokenSource.Cancel();
        }
    }
}