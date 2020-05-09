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
            [Perper("queries")] IAsyncDisposable queries,
            [Perper("gossips")] IAsyncDisposable gossips,
            [PerperStream("output")] IAsyncCollector<IAsyncDisposable> output,
            CancellationToken cancellationToken)
        {
            foreach (var node in nodes)
            {
                var proposer = await context.StreamFunctionAsync(typeof(Proposer), new {node, nodes, queries});
                var validator = await context.StreamFunctionAsync(typeof(Validator), new {node, queries});
                var acceptor = await context.StreamFunctionAsync(typeof(Acceptor), new {node, nodes, gossips, proposer, validator});
                var committer = await context.StreamFunctionAsync(typeof(Committer), new {node, acceptor});

                await Task.WhenAll(new[] {proposer, validator, acceptor, committer}.Select(
                    stream => output.AddAsync(stream, cancellationToken)));
            }
        }
    }
}