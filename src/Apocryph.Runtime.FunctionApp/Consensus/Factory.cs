using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Runtime.FunctionApp.Consensus.Core;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Runtime.FunctionApp.Consensus
{
    public class Factory
    {
        [FunctionName(nameof(Factory))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("nodes")] Node[] nodes,
            [Perper("localNodes")] IAsyncEnumerable<Node> localNodes,
            [Perper("ibc")] IAsyncDisposable ibc,
            [Perper("queries")] IAsyncDisposable queries,
            [Perper("gossips")] IAsyncDisposable gossips,
            [PerperStream("output")] IAsyncCollector<IAsyncDisposable> output,
            CancellationToken cancellationToken)
        {
            // TODO: Have a way to remove local nodes
            await foreach (var node in localNodes)
            {
                var proposer = await context.StreamFunctionAsync(typeof(Proposer), new { node, nodes, queries, ibc });
                var validator = await context.StreamFunctionAsync(typeof(Validator), new { node, queries, ibc });
                var committer = await context.StreamFunctionAsync(typeof(Committer), new { node, nodes, gossips, proposer, validator });

                await Task.WhenAll(new[] { proposer, validator, committer }.Select(
                    stream => output.AddAsync(stream, cancellationToken)));
            }
        }
    }
}