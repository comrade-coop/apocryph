using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Model;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Bindings;
using Perper.WebJobs.Extensions.Model;
using Perper.WebJobs.Extensions.Triggers;

namespace Apocryph.FunctionApp
{
    public static class _Fixup
    {
        public static async Task Listen<T>(this IAsyncEnumerable<T> stream,
            Action<T> function, CancellationToken ct)
        {
            await foreach(var x in stream) {
                function.Invoke(x);
            }
        }

        public static async Task Listen<T>(this IAsyncEnumerable<T> stream,
            Func<T, Task> function, CancellationToken ct)
        {
            await foreach(var x in stream) {
                await function.Invoke(x);
            }
        }

        public static Task SetState<T>(this IPerperStreamContext stream,
            string name, T state)
        {
            throw new NotImplementedException();
        }

        public static Task<T> GetState<T>(this IPerperStreamContext stream,
            string name)
        {
            throw new NotImplementedException();
        }
    }
}