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
            [Perper("slotGossips")] IPerperStream slotGossips,
            [Perper("chains")] IDictionary<Guid, Chain> chains,
            [Perper("output")] IAsyncCollector<int> output,
            CancellationToken cancellationToken)
        {
            await using var gossips = context.DeclareStream("Peering-gossips", typeof(PeeringStream));
            await using var queries = context.DeclareStream("Peering-queries", typeof(PeeringStream));
            await using var reports = context.DeclareStream("Peering-reports", typeof(PeeringStream));
            await using var salts = context.DeclareStream("Salts", typeof(SaltsStream));

            await using var chain = await context.StreamFunctionAsync("Chain", typeof(ChainStream), new
            {
                chains,
                gossips,
                queries,
                salts = salts.Subscribe(),
                slotGossips = slotGossips.Subscribe()
            });

            await using var validator = await context.StreamFunctionAsync("DummyStream", new
            {
                queries = queries.Subscribe(), // HACK: Make sure the queries peering receives all streams
                gossips = gossips.Subscribe(), // HACK: Make sure the gossips peering receives all streams
                reports = reports.Subscribe(),
            });

            await using var ibc = await context.StreamFunctionAsync("IBC-global", typeof(IBCStream), new
            {
                chain = chain.Subscribe(),
                validator = validator.Subscribe(),
                gossips = gossips.Subscribe(),
                nodes = new Dictionary<Guid, Node?[]>(),
                node = default(Node?)
            });
            await using var filter = await context.StreamFunctionAsync("Filter-global", typeof(FilterStream), new
            {
                ibc = ibc.Subscribe(),
                gossips = gossips.Subscribe(),
                chains
            });

            await context.StreamFunctionAsync(salts, new
            {
                chains,
                filter = filter.Subscribe()
            });

            await context.StreamFunctionAsync(gossips, new
            {
                factory = chain.Subscribe(),
                filter = typeof(IBCStream)
            });

            await context.StreamFunctionAsync(queries, new
            {
                factory = chain.Subscribe(),
                filter = typeof(ConsensusStream)
            });

            await context.StreamFunctionAsync(reports, new
            {
                factory = chain.Subscribe(),
                filter = typeof(ConsensusStream)
            });

            await using var loggingStream = await context.StreamActionAsync(typeof(LoggingStream), new
            {
                chain = chain.Subscribe(),
                nodes = new Dictionary<Guid, Node?[]>(),
                filter = filter.Subscribe(),
                reports = reports.Subscribe()
            });

            await context.BindOutput(cancellationToken);
        }
    }
}