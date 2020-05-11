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
    public class Acceptor
    {
        [FunctionName(nameof(Acceptor))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("node")] Node node,
            [PerperStream("committer")] IAsyncEnumerable<Message<Block>> committer,
            [PerperStream("output")] IAsyncCollector<Block> output,
            CancellationToken cancellationToken)
        {
            var acceptedMessage = await committer.FirstAsync(cancellationToken);
            await output.AddAsync(acceptedMessage.Value, cancellationToken);
        }
    }
}