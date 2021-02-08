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
            Assert.Equal(service.Inputs, result!.Inputs, new DictionaryComparer<string, string>());
            Assert.Equal(service.Outputs, result!.Outputs, new DictionaryComparer<string, IStream>());
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
            Assert.Equal(service!.Inputs, result!.Inputs, new DictionaryComparer<string, string>());
            Assert.Equal(service!.Outputs, result!.Outputs, new DictionaryComparer<string, IStream>());
        }

        class DictionaryComparer<TK, TV> : IEqualityComparer<IDictionary<TK, TV>> where TK : notnull
        {
            public IEqualityComparer<TV> ValueComparer { get; set; } = EqualityComparer<TV>.Default;

            public bool Equals(IDictionary<TK, TV> a, IDictionary<TK, TV> b)
            {
                if (a.Count != b.Count) return false;

                foreach (var (key, valueA) in a)
                {
                    if (!b.TryGetValue(key, out var valueB) || !ValueComparer.Equals(valueA, valueB))
                    {
                        return false;
                    }
                }

                return true;
            }

            public int GetHashCode(IDictionary<TK, TV> x)
            {
                throw new NotImplementedException();
            }
        }
    }
}