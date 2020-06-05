using System.Collections.Generic;
using System.Threading;
using System.Threading.Tasks;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Runtime.FunctionApp
{
    public class ChainList
    {
        [FunctionName(nameof(ChainList))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("chains")] IDictionary<byte[], string> chains,
            CancellationToken cancellationToken)
        {
            var gossips = context.DeclareStream(typeof(Peering));
            var queries = context.DeclareStream(typeof(Peering));

            var miner = await context.StreamFunctionAsync(typeof(Miner), new { chains });

            var chain = await context.StreamFunctionAsync(typeof(Chain), new { miner, gossips, queries });
            await context.StreamFunctionAsync(gossips, new
            {
                chain, filter = new[]
                {
                    typeof(Assigner),
                    typeof(Proposer),
                    typeof(Committer)
                }
            });
            await context.StreamFunctionAsync(queries, new
            {
                chain, filter = new[]
                {
                    typeof(Proposer)
                }
            });
        }
    }
}