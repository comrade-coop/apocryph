using Apocryph.PerperUtilities;
using Perper.WebJobs.Extensions.Fake;

namespace Apocryph.ServiceRegistry.Test
{
    using ServiceRegistry = Apocryph.ServiceRegistry.FunctionApp.ServiceRegistry;

    public static class ServiceRegistryFakes
    {
        public static (FakeAgent, ServiceRegistry) GetServiceRegistryAgent()
        {
            var registry = new ServiceRegistry(new FakeState());

            var agent = new FakeAgent();
            agent.RegisterFunction("Register", ((ServiceLocator locator, Service service) input) => registry.Register(input, default));
            agent.RegisterFunction("RegisterTypeHandler", ((string type, Handler handler) input) => registry.RegisterTypeHandler(input, default));
            agent.RegisterFunction("Lookup", (ServiceLocator input) => registry.Lookup(input, default));

            return (agent, registry);
        }
    }
}