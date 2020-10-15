using System.Collections.Generic;
using System.Threading.Tasks;
using System.Text.Json;
using Block = Apocryph.Core.Consensus.Blocks.Block;
using Apocryph.Core.Consensus.Serialization;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;
using Ipfs.Http;

namespace Apocryph.Runtime.FunctionApp
{
    public class HashRegistryStream
    {
        [FunctionName(nameof(HashRegistryStream))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("input")] IAsyncEnumerable<Block> input)
        {
            var ipfs = new IpfsClient();

            await foreach (var value in input)
            {
                var serialized = JsonSerializer.SerializeToUtf8Bytes(value, ApocryphSerializationOptions.JsonSerializerOptions);

                // FIXME: Using "raw" here instead of "json", since Ipfs.Http.Client doesn't seem to consider "json" a valid MultiCodec
                await ipfs.Block.PutAsync(serialized, "raw", "sha2-256");
            }
        }
    }
}