using System.Collections.Generic;
using System.Threading.Tasks;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Peering
{
    [PerperData]
    public class PeerHandler
    {
        public IAgent Agent { get; private set; }

        public string Method { get; private set; }

        public PeerHandler(IAgent agent, string method)
        {
            Agent = agent;
            Method = method;
        }

        public virtual Task<IAsyncEnumerable<object>> InvokeAsync(object? parameters)
        {
            return Agent.CallFunctionAsync<IAsyncEnumerable<object>>(Method, parameters);
        }
    }
}