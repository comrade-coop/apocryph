using System;
using System.Collections.Generic;
using System.Threading;
using System.Threading.Tasks;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;
using Apocryph.Core.Consensus.Blocks;
using Apocryph.Core.Consensus.VirtualNodes;

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
            var salts = context.DeclareStream(typeof(SaltsStream));

            var chain = await context.StreamFunctionAsync(typeof(ChainStream), new { chains, gossips, queries, salts, slotGossips });
            await output.AddAsync(chain);

            var node = new Node(new byte[0], -1);
            var ibc = await context.StreamFunctionAsync(typeof(IBCStream), new { chain, gossips, node, nodes = new Dictionary<byte[], Node?[]>() });
            var filter = await context.StreamFunctionAsync(typeof(FilterStream), new { ibc, gossips, chains, node });

            await context.StreamFunctionAsync(salts, new { chains, filter });

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