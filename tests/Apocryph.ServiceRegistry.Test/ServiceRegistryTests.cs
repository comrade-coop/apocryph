using System.Collections.Generic;
using Apocryph.PerperUtilities;
using Perper.WebJobs.Extensions.Fake;
using Perper.WebJobs.Extensions.Model;
using Xunit;

namespace Apocryph.ServiceRegistry.Test
{
    using ServiceRegistry = Apocryph.ServiceRegistry.FunctionApp.ServiceRegistry;

    public class ServiceRegistryTests
    {
        private Service GetTestService()
        {
            var stream = new FakeStream<int>(new[] { 0, 1, 2 });
            return new Service(new Dictionary<string, string>() {
                {"numbers", "numbers"}
            }, new Dictionary<string, IStream>() {
                {"numbers", stream}
            });
        }

        [Fact]
        public async void Lookup_BasicRegister_ReturnsExact()
        {
            var registry = new ServiceRegistry(new FakeState());

            var locator = new ServiceLocator("number", "0");
            var service = GetTestService();

            await registry.Register((locator, service), default);

            var result = await registry.Lookup(locator, default);

            Assert.NotNull(result);
            Assert.Equal(service.Inputs, result!.Inputs);
            Assert.Equal(service.Outputs, result!.Outputs);
        }

        [Fact]
        public async void Lookup_Unregistered_ReturnsNull()
        {
            var registry = new ServiceRegistry(new FakeState());

            var locator = new ServiceLocator("number", "0");
            var service = GetTestService();

            // A - Before registering
            Assert.Null(await registry.Lookup(locator, default));

            await registry.Register((locator, service), default);

            // B - After registering
            var wrongLocator1 = new ServiceLocator("number", "1");
            Assert.Null(await registry.Lookup(wrongLocator1, default));

            var wrongLocator2 = new ServiceLocator("string", "0");
            Assert.Null(await registry.Lookup(wrongLocator2, default));
        }

        [Fact]
        public async void Lookup_WithHandler_CallsHandlerAndReturns()
        {
            var registry = new ServiceRegistry(new FakeState());
            var locator = new ServiceLocator("number", "0");
            Service? service = null;

            var agent = new FakeAgent();
            agent.RegisterFunction("method", async (object? parameters) =>
            {
                service = GetTestService();
                await registry.Register((locator, service), default);
            });
            var handler = new Handler(agent, "method");

            await registry.RegisterTypeHandler((locator.Type, handler), default);

            Assert.Null(service);

            var result = await registry.Lookup(locator, default);
            Assert.NotNull(result);
            Assert.NotNull(service);
            Assert.Equal(service!.Inputs, result!.Inputs);
            Assert.Equal(service!.Outputs, result!.Outputs);
        }
    }
}