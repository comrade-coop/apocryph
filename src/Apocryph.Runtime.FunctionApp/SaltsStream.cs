using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Core.Consensus;
using Apocryph.Core.Consensus.Blocks;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Runtime.FunctionApp
{
    public class SaltsStream
    {
        [FunctionName(nameof(SaltsStream))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("chains")] Dictionary<Guid, Chain> chains,
            [Perper("filter")] IAsyncEnumerable<Hash> filter,
            [Perper("hashRegistryWorker")] string hashRegistryWorker,
            [Perper("output")] IAsyncCollector<(Guid, int, byte[])> output,
            CancellationToken cancellationToken)
        {
            await foreach (var hash in filter)
            {
                var block = await context.CallWorkerAsync<Block>(hashRegistryWorker, new { hash }, cancellationToken);
                var chain = chains[block!.ChainId];
                foreach (var (slot, salt) in RandomWalk.Run(hash).Take(1 + chain.SlotCount / 10))
                {
                    await output.AddAsync((block!.ChainId, (int)(slot % chain.SlotCount), salt.Value));
                }
            }
        }
    }
}