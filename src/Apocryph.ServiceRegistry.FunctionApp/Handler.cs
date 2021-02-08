using System.Threading.Tasks;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.ServiceRegistry
{
    [PerperData]
    public class Handler
    {
        public Handler(IAgent agent, string method)
        {
            Agent = agent;
            Method = method;
        }

        public IAgent Agent { get; private set; }
        public string Method { get; private set; }

        public Task InvokeAsync(object? parameters)
        {
            return Agent.CallActionAsync(Method, parameters);
        }
    }
}