using System.Threading.Tasks;
using Perper.WebJobs.Extensions.Fake;
using Xunit;

namespace Apocryph.HashRegistry.Test
{
    using HashRegistry = Apocryph.HashRegistry.FunctionApp.HashRegistry;

    public class HashRegistryTests
    {
        public static object[][] SampleData = new[] {
            new object[] { new byte[] { 0, 1, 2, 255, 254, 253 } }
        };

        [Theory]
        [MemberData(nameof(SampleData))]
        public async void Retrieve_AfterStore_ReturnsExact(byte[] data)
        {
            var registry = new HashRegistry(new FakeState());
            var hash = Hash.FromSerialized<object?>(data);

            await registry.Store(data, default);

            var result = await registry.Retrieve(hash, default);

            Assert.Equal(data, result, new ArrayComparer<byte>());
        }

        [Theory]
        [MemberData(nameof(SampleData))]
        public async void Retrieve_BeforeStore_ReturnsAfter(byte[] data)
        {
            const int waitTime = 100;

            var registry = new HashRegistry(new FakeState());
            var hash = Hash.FromSerialized<object?>(data);

            var resultTask = registry.Retrieve(hash, default);

            await Task.Delay(waitTime);
            Assert.False(resultTask.IsCompleted);

            await registry.Store(data, default);

            var result = await resultTask;

            Assert.Equal(data, result, new ArrayComparer<byte>());
        }
    }
}