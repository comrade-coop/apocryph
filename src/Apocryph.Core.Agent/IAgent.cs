using System;
using System.Threading.Tasks;

namespace Apocryph.Core.Agent
{
    public interface IAgent<T> where T : class
    {
        void Setup(IContext<T> context);
        Task Run(IContext<T> context, object message, Guid? reference);
    }
}