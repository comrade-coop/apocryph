using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Runtime.FunctionApp.Consensus.Core;
using Apocryph.Runtime.FunctionApp.ValidatorSelection;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Runtime.FunctionApp
{
    public class ChainFactory
    {
        [FunctionName(nameof(ChainFactory))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("chains")] IAsyncEnumerable<object> chains,
            [Perper("miner")] IAsyncDisposable miner,
            [Perper("gossips")] IAsyncDisposable gossips,
            [Perper("ibc")] IAsyncDisposable ibc,
            [PerperStream("output")] IAsyncCollector<IAsyncDisposable> output,
            CancellationToken cancellationToken)
        {
            await foreach (var chain in chains)
            {
                var nodes = new Node[] { new Node { Id = 1 } }; // chain.Nodes
                var chainId = new byte[] { 0 }; // chain.Id
                IAsyncDisposable salts = default!;

                var localNodes = await context.StreamFunctionAsync(typeof(Assigner), new { nodes, chainId, gossips, salts, miner });
                var chainStream = await context.StreamFunctionAsync(typeof(Chain), new { nodes, localNodes, ibc });

                await Task.WhenAll(new[] { localNodes, chainStream }.Select(
                    stream => output.AddAsync(stream, cancellationToken)));
            }
        }
    }
}