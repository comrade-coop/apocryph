using System.Threading.Tasks;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.PerperUtilities
{
    [PerperData]
    public class Handler : IHandler
    {
        public Handler(IAgent agent, string method)
        {
            Agent = agent;
            Method = method;
        }

        public IAgent Agent { get; private set; }
        public string Method { get; private set; }

        public virtual Task InvokeAsync(object? parameters)
        {
            return Agent.CallActionAsync(Method, parameters);
        }
    }

    [PerperData]
    public class Handler<T> : Handler, IHandler<T>
    {
        public Handler(IAgent agent, string method)
            : base(agent, method) { }

        public new virtual Task<T> InvokeAsync(object? parameters)
        {
            return Agent.CallFunctionAsync<T>(Method, parameters);
        }
    }
}