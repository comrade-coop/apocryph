using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Runtime.FunctionApp
{
    public class PeeringStream
    {
        [FunctionName(nameof(PeeringStream))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("factory")] IAsyncEnumerable<IPerperStream> factory,
            [Perper("filter")] Type filter,
            CancellationToken cancellationToken)
        {
            var output = new List<IPerperStream>();
            var peers = factory; //.Where(stream => stream.GetDelegate() == filter);
            await foreach (var peer in peers.WithCancellation(cancellationToken))
            {
                output.Add(peer);
                await context.RebindOutput(output);
            }
        }
    }
}