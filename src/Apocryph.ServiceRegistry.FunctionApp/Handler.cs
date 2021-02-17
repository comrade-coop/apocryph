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

        public virtual Task InvokeAsync(object? parameters)
        {
            return Agent.CallActionAsync(Method, parameters);
        }
    }

    [PerperData]
    public class HandlerWithContext : Handler
    {
        public HandlerWithContext(IAgent agent, string method, object? context)
            : base(agent, method)
        {
            Context = context;
        }

        public object? Context { get; private set; }

        public override Task InvokeAsync(object? parameters)
        {
            return Agent.CallActionAsync(Method, (Context, parameters));
        }
    }
}