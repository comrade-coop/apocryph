using System;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Runtime.FunctionApp.Consensus.Core;
using Apocryph.Runtime.FunctionApp.Consensus;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Runtime.FunctionApp
{
    public class Chain
    {
        [FunctionName(nameof(Chain))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("nodes")] Node[] nodes,
            [PerperStream("localNodes")] IAsyncDisposable localNodes,
            [PerperStream("ibc")] IAsyncDisposable ibc,
            CancellationToken cancellationToken)
        {
            var queries = context.DeclareStream(typeof(Peering));
            var gossips = context.DeclareStream(typeof(Peering));

            var factory = await context.StreamFunctionAsync(typeof(Factory), new { nodes, localNodes, ibc, queries, gossips });
            await context.StreamFunctionAsync(queries, new { factory, filter = typeof(Proposer) });
            await context.StreamFunctionAsync(gossips, new { factory, filter = typeof(Committer) });

            var output = await context.StreamFunctionAsync(typeof(Acceptor), new { nodes, gossips });

            await context.BindOutput(output, cancellationToken);
        }
    }
}