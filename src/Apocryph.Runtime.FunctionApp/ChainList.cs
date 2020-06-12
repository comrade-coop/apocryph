using System;
using System.Collections.Generic;
using System.Threading;
using System.Threading.Tasks;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Runtime.FunctionApp
{
    public class ChainList
    {
        [FunctionName(nameof(ChainList))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("slotGossips")] IAsyncDisposable slotGossips,
            [Perper("chains")] IDictionary<byte[], int> chains,
            [PerperStream("output")] IAsyncCollector<IAsyncDisposable> output,
            CancellationToken cancellationToken)
        {
            var gossips = context.DeclareStream(typeof(Peering));
            var queries = context.DeclareStream(typeof(Peering));
            var salts = (IAsyncDisposable)default!;

            var chainStreams = new List<IAsyncDisposable>();
            foreach (var (chainId, slotCount) in chains)
            {
                var chain = await context.StreamFunctionAsync(typeof(Chain), new { chainId, slotCount, gossips, queries, salts, slotGossips });
                chainStreams.Add(chain);

                await output.AddAsync(chain);
            }

            await context.StreamFunctionAsync(gossips, new
            {
                factory = chainStreams.ToArray(),
                filter = new[]
                {
                    typeof(IBC)
                }
            });

            await context.StreamFunctionAsync(queries, new
            {
                factory = chainStreams.ToArray(),
                filter = new[]
                {
                    typeof(Consensus)
                }
            });
        }
    }
}