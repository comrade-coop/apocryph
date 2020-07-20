using System;
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
            [Perper("chains")] Dictionary<Guid, Chain> chains,
            [Perper("filter")] IAsyncEnumerable<Block> filter,
            [Perper("output")] IAsyncCollector<(Guid, int, byte[])> output,
            CancellationToken cancellationToken)
        {
            var options = new JsonSerializerOptions
            {
                Converters =
                {
                    { new NonStringKeyDictionaryConverter() }
                }
            };

            await foreach (var block in filter)
            {
                var chain = chains[block.ChainId];
                foreach (var (slot, salt) in RandomWalk.Run(JsonSerializer.SerializeToUtf8Bytes(block, options)).Take(1 + chain.SlotCount / 10))
                {
                    await output.AddAsync((block.ChainId, (int)(slot % chain.SlotCount), salt));
                }
            }
        }
    }
}