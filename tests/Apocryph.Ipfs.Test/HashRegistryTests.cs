using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Xunit;

namespace Apocryph.Ipfs.Test
{
    [Collection("Ipfs Collection")]
    public class HashResolverTests
    {
        private IpfsFixture _fixture;

        public HashResolverTests(IpfsFixture fixture)
        {
            _fixture = fixture;
        }

        public static IEnumerable<object?[]> SampleData() =>
            IpfsFixture.HashResolverImplementations.SelectMany((i) =>
                TestData.DataInterface.Select((d) => new object?[] { i, d[0] }));

        [Theory]
        [MemberData(nameof(SampleData))]
        public async void Retrieve_AfterStore_ReturnsExact(string resolverImplementation, IExample data)
        {
            var resolver = _fixture.GetHashResolver(resolverImplementation, 1);

            var hash = Hash.From(data);

            var storedHash = await resolver.StoreAsync(data);

            Assert.Equal(hash, storedHash);

            var result = await resolver.RetrieveAsync(hash);

            Assert.Equal(data, result);
        }

        [Theory]
        [MemberData(nameof(SampleData))]
        public async void Retrieve_BeforeStore_ReturnsAfter(string resolverImplementation, IExample _data)
        {
            var resolver = _fixture.GetHashResolver(resolverImplementation, 2);

            var data = ("testwrap", _data);

            var hash = Hash.From(data);

            var resultTask = resolver.RetrieveAsync(hash);

            await Task.Delay(TestData.WaitTime);
            Assert.False(resultTask.IsCompleted);

            await resolver.StoreAsync(data);

            var result = await resultTask;

            Assert.Equal(data, result);
        }

        [Theory]
        [MemberData(nameof(SampleData))]
        public async void Retrieve_AfterRemoteStore_ReturnsExact(string resolverImplementation, IExample data)
        {
            var resolverFrom = _fixture.GetHashResolver(resolverImplementation, 1);
            var resolverTo = _fixture.GetHashResolver(resolverImplementation, 2);

            var hash = Hash.From(data);

            var storedHash = await resolverFrom.StoreAsync(data);

            Assert.Equal(hash, storedHash);

            var cts = new CancellationTokenSource(TimeSpan.FromSeconds(1));
            var result = await resolverTo.RetrieveAsync(hash, cts.Token);

            Assert.Equal(data, result);
        }
    }
}