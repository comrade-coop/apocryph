using System.Collections.Generic;
using System.Linq;
using System.Text.Json;
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
            [Perper("chains")] Dictionary<byte[], Chain> chains,
            [Perper("filter")] IAsyncEnumerable<Block> filter,
            [Perper("output")] IAsyncCollector<(byte[], int, byte[])> output,
            CancellationToken cancellationToken)
        {
            await foreach (var block in filter)
            {
                var chain = chains[block.ChainId];
                foreach (var (slot, salt) in RandomWalk.Run(JsonSerializer.SerializeToUtf8Bytes(block)).Take(chain.SlotCount / 10))
                {
                    await output.AddAsync((block.ChainId, (int)(slot % chain.SlotCount), salt));
                }
            }
        }
    }
}