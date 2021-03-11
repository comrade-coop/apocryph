using Perper.WebJobs.Extensions.Model;

namespace Apocryph.PerperUtilities
{
    public static class HandlerExtensions
    {
        public static IHandler GetHandler(this IAgent agent, string methodName) => new Handler(agent, methodName);
        public static IHandler<T> GetHandler<T>(this IAgent agent, string methodName) => new Handler<T>(agent, methodName);
    }
}