using System.Threading;
using System.Threading.Tasks;
using Apocryph.Runtime.FunctionApp.Consensus.Core;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Runtime.FunctionApp.Consensus
{
    public class Consensus
    {
        [FunctionName(nameof(Consensus))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("nodes")] Node[] nodes,
            CancellationToken cancellationToken)
        {
            var queries = context.DeclareStream(typeof(Peering));
            var gossips = context.DeclareStream(typeof(Peering));
            var output = context.DeclareStream(typeof(Peering));

            var factory = await context.StreamFunctionAsync(typeof(Factory), new {nodes, queries, gossips});
            await context.StreamFunctionAsync(queries, new {factory, filter = typeof(Proposer)});
            await context.StreamFunctionAsync(gossips, new {factory, filter = typeof(Acceptor)});
            await context.StreamFunctionAsync(output, new {factory, filter = typeof(Committer)});

            await context.BindOutput(output, cancellationToken);
        }
    }
}