using Xunit;

namespace Apocryph.HashRegistry.Test
{
    public class HashRegistryProxyTests
    {
        class Example
        {
            public int Number { get; set; }
        }

        [Theory]
        [InlineData(10)]
        public async void Retrieve_AfterStore_ReturnsExact(int data)
        {
            var proxy = HashRegistryFakes.GetHashRegistryProxy();

            var value = new Example { Number = data };

            var hash = await proxy.StoreAsync(value, default);

            var result = await proxy.RetrieveAsync<Example>(hash, default);

            Assert.Equal(data, result.Number);
            Assert.NotEqual(value, result); // Reference equality
        }
    }
}