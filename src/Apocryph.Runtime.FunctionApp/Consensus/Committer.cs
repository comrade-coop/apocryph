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
    public class Committer
    {
        [FunctionName(nameof(Committer))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("node")] Node node,
            [PerperStream("acceptor")] IAsyncEnumerable<Gossip<Block>> acceptor,
            [PerperStream("output")] IAsyncCollector<Block> output,
            CancellationToken cancellationToken)
        {
            var committedGossip = await acceptor.FirstAsync(
                gossip => gossip.Signers.Contains(node) && gossip.Verb == GossipVerb.Commit,
                cancellationToken);
            await output.AddAsync(committedGossip.Value, cancellationToken);
        }
    }
}