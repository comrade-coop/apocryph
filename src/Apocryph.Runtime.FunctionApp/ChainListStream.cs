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
            await using var peering = context.DeclareStream("Peering", typeof(PeeringStream));
            var gossips = peering;
            var queries = peering;
            var reports = peering;
            var hashes = peering;
            await using var salts = context.DeclareStream("Salts", typeof(SaltsStream));

            await using var hashRegistry = await context.StreamFunctionAsync("HashRegistry", typeof(HashRegistryStream), new
            {
                filter = typeof(Block),
                input = hashes.Subscribe()
            }, typeof(HashRegistryEntry));

            await context.StreamActionAsync("DummyStream", new
            {
                hashRegistry = hashRegistry.Subscribe() // HACK: Ensure hash registry is started up before anything else
            });

            await using var chain = await context.StreamFunctionAsync("Chain", typeof(ChainStream), new
            {
                chains,
                gossips,
                queries,
                hashRegistry,
                salts = salts.Subscribe(),
                slotGossips = slotGossips.Subscribe()
            });

            // HACK: Create an empty stream for the global IBC
            await using var validator = await context.StreamFunctionAsync("DummyStream", new
            {
                peering = peering.Subscribe(), // HACK: Ensure peering is started before it starts receiving streams
                hashRegistry = hashRegistry.Subscribe() // HACK: Ensure hash registry is started up
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
                hashRegistry,
                chains
            });

            await context.StreamFunctionAsync(salts, new
            {
                chains,
                hashRegistry = hashRegistry,
                filter = filter.Subscribe()
            });

            await context.StreamFunctionAsync(peering, new
            {
                factory = chain.Subscribe(),
                initial = new List<IPerperStream>() { filter },
            });

            await using var reportsStream = await context.StreamActionAsync(typeof(ReportsStream), new
            {
                hashRegistry = hashRegistry,
                chain = chain.Subscribe(),
                nodes = new Dictionary<Guid, Node?[]>(),
                filter = filter.Subscribe(),
                reports = reports.Subscribe()
            });

            await context.BindOutput(cancellationToken);
        }
    }
}