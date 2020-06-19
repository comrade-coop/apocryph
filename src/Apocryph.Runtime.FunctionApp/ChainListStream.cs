using System;
using System.Collections.Generic;
using System.Threading;
using System.Threading.Tasks;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;
using Apocryph.Core.Consensus.Blocks;

namespace Apocryph.Runtime.FunctionApp
{
    public class ChainListStream
    {
        [FunctionName(nameof(ChainListStream))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("slotGossips")] IAsyncDisposable slotGossips,
            [Perper("chains")] IDictionary<byte[], Chain> chains,
            [PerperStream("output")] IAsyncCollector<IAsyncDisposable> output,
            CancellationToken cancellationToken)
        {
            var gossips = context.DeclareStream(typeof(PeeringStream));
            var queries = context.DeclareStream(typeof(PeeringStream));
            var salts = (IAsyncDisposable)default!;

            var chain = await context.StreamFunctionAsync(typeof(ChainStream), new { chains, gossips, queries, salts, slotGossips });
            await output.AddAsync(chain);

            await context.StreamFunctionAsync(gossips, new
            {
                factory = chain,
                filter = new[]
                {
                    typeof(IBCStream)
                }
            });

            await context.StreamFunctionAsync(queries, new
            {
                factory = chain,
                filter = new[]
                {
                    typeof(ConsensusStream)
                }
            });
        }
    }
}