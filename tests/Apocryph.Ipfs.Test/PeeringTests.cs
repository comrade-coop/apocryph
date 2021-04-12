using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
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

        public static IEnumerable<object?[]> SampleData() =>
            IpfsFixture.PeerConnectorImplementations.SelectMany((i) =>
                TestData.DataInterface.Select((d) => new object?[] { i, d[0], d[0] }));

        [Theory]
        [MemberData(nameof(SampleData))]
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
                // Assert.NotEqual(request, dataRequest, ReferenceEqualityComparer.Instance); // NOTE .NET 5

                return dataRequest;
            }, cancellationTokenSource.Token);

            var reply = await connectorFrom.Query<IExample, IExample>(await connectorTo.Self, path, dataRequest, cancellationTokenSource.Token);
            Assert.Equal(reply, dataReply, SerializedComparer.Instance);
            // Assert.NotEqual(reply, dataReply, ReferenceEqualityComparer.Instance); // NOTE .NET 5
            cancellationTokenSource.Cancel();
        }
    }
}