using System.Threading;
using System.Threading.Tasks;
using Apocryph.PerperUtilities;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Model;
using Perper.WebJobs.Extensions.Triggers;

namespace Apocryph.ServiceRegistry.FunctionApp
{
    public class ServiceRegistry
    {
        private IState _state;

        public ServiceRegistry(IState state)
        {
            _state = state;
        }

        public string GetLocatorKey(ServiceLocator locator) => $"service:{locator.Type}:{locator.Id}";
        public string GetHandlerKey(string type) => $"handler:{type}";

        [FunctionName("Register")]
        public async Task Register([PerperTrigger] (ServiceLocator locator, Service service) parameters, CancellationToken cancellationToken)
        {
            await _state.SetValue(GetLocatorKey(parameters.locator), parameters.service);
        }

        [FunctionName("RegisterTypeHandler")]
        public async Task RegisterTypeHandler([PerperTrigger] (string type, Handler handler) parameters, CancellationToken cancellationToken)
        {
            await _state.SetValue(GetHandlerKey(parameters.type), parameters.handler);
        }

        [FunctionName("Lookup")]
        public async Task<Service?> Lookup([PerperTrigger] ServiceLocator locator, CancellationToken cancellationToken)
        {
            var service = await _state.GetValue<Service?>(GetLocatorKey(locator), () => null);
            if (service == null)
            {
                var handler = await _state.GetValue<Handler?>(GetHandlerKey(locator.Type), () => null);
                if (handler != null)
                {
                    await handler.InvokeAsync(locator);
                    service = await _state.GetValue<Service?>(GetLocatorKey(locator), () => null);
                }
            }
            return service;
        }
    }
}