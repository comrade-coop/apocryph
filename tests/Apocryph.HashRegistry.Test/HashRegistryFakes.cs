using Perper.WebJobs.Extensions.Fake;

namespace Apocryph.HashRegistry.Test
{
    using HashRegistry = Apocryph.HashRegistry.FunctionApp.HashRegistry;

    public static class HashRegistryFakes
    {
        public static HashRegistryProxy GetHashRegistryProxy()
        {
            var registry = new HashRegistry(new FakeState());

            var agent = new FakeAgent();
            agent.RegisterFunction("Store", (byte[] data) => registry.Store(data, default));
            agent.RegisterFunction("Retrieve", (Hash hash) => registry.Retrieve(hash, default));

            return new HashRegistryProxy(agent);
        }
    }
}