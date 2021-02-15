using Perper.WebJobs.Extensions.Fake;
using Xunit;

namespace Apocryph.HashRegistry.Test
{
    using HashRegistry = Apocryph.HashRegistry.FunctionApp.HashRegistry;

    public class HashRegistryProxyTests
    {
        private HashRegistryProxy GetHashRegistryProxy()
        {
            var registry = new HashRegistry(new FakeState());

            var agent = new FakeAgent();
            agent.RegisterFunction("Store", (byte[] data) => registry.Store(data, default));
            agent.RegisterFunction("Retrieve", (Hash hash) => registry.Retrieve(hash, default));

            return new HashRegistryProxy(agent);
        }

        class Example
        {
            public int Number { get; set; }
        }

        [Theory]
        [InlineData(10)]
        public async void Retrieve_AfterStore_ReturnsExact(int data)
        {
            var proxy = GetHashRegistryProxy();

            var value = new Example { Number = data };

            var hash = await proxy.StoreAsync(value, default);

            var result = await proxy.RetrieveAsync<Example>(hash, default);

            Assert.Equal(data, result.Number);
            Assert.NotEqual(value, result); // Reference equality
        }
    }
}