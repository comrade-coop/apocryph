using System.Threading;
using System.Threading.Tasks;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.FunctionApp
{
    public static class _Fixup
    {
        public static Task CallStreamAction(this IPerperStreamContext context, string name, object parameters, CancellationToken cancellationToken)
        {
            return context.CallStreamAction(name, parameters);
        }
    }
}