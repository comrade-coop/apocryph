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
            [Perper("slotGossips")] IPerperStream slotGossips,
            [Perper("chains")] IDictionary<byte[], Chain> chains,
            [Perper("output")] IAsyncCollector<IPerperStream> output,
            CancellationToken cancellationToken)
        {
            await using var gossips = context.DeclareStream(typeof(PeeringStream));
            await using var queries = context.DeclareStream(typeof(PeeringStream));
            await using var salts = context.DeclareStream(typeof(SaltsStream));

            var chain = await context.StreamFunctionAsync(typeof(ChainStream), new
            {
                chains,
                gossips,
                queries,
                salts,
                slotGossips
            });
            await output.AddAsync(chain);

            var node = new Node(new byte[0], -1);
            var ibc = await context.StreamFunctionAsync(typeof(IBCStream), new
            {
                chain = chain.Subscribe(),
                gossips = gossips.Subscribe(),
                node,
                nodes = new Dictionary<byte[], Node?[]>()
            });
            var filter = await context.StreamFunctionAsync(typeof(FilterStream), new
            {
                ibc = ibc.Subscribe(),
                gossips = gossips.Subscribe(),
                chains,
                node
            });

            await context.StreamActionAsync(salts, new
            {
                chains,
                filter = filter.Subscribe()
            });

            await context.StreamActionAsync(gossips, new
            {
                factory = chain.Subscribe(),
                filter = typeof(IBCStream)
            });

            await context.StreamActionAsync(queries, new
            {
                factory = chain.Subscribe(),
                filter = typeof(ConsensusStream)
            });

            await context.BindOutput(cancellationToken);
        }
    }
}