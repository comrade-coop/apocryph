using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Model;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.FunctionApp
{
    public static class StreamJoiner
    {
        [FunctionName("StreamJoiner")]
        public static async Task Run([PerperStreamTrigger("StreamJoiner")] IPerperStreamContext context,
            [Perper("streamA")] IAsyncEnumerable<object> streamA,
            [Perper("streamB")] IAsyncEnumerable<object> streamB,
            [Perper("outputStream")] IAsyncCollector<object> outputStream)
        {
            await Task.WhenAll(
                streamA.ForEachAsync(x => outputStream.AddAsync(x), CancellationToken.None),
                streamB.ForEachAsync(x => outputStream.AddAsync(x), CancellationToken.None));
        }

        public static async Task<IAsyncDisposable> JoinStreams(this IPerperStreamContext context,
            params IAsyncDisposable[] streams)
        {
            // NOTE: Should probably dispose all the intermediate streams as well
            // NOTE: Could do with O(log(N)) routing complexity instead of O(N)

            var lastStream = streams.First();

            foreach (var stream in streams.Skip(1))
            {
                lastStream = await context.StreamFunctionAsync("StreamJoiner", new
                {
                    streamA = lastStream,
                    streamB = stream
                });
            }

            return lastStream;
        }
    }
}