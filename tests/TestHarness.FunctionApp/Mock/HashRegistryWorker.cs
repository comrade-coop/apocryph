using System;
using System.Threading;
using System.Threading.Tasks;
using System.Text.Json;
using Apocryph.Core.Consensus.Blocks;
using Apocryph.Core.Consensus.Serialization;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace TestHarness.FunctionApp.Mock
{
    public class HashRegistryWorker
    {
        [FunctionName(nameof(HashRegistryWorker))]
        [return: Perper("$return")]
        public async Task<Block> Run([PerperWorkerTrigger] PerperWorkerContext context,
            [Perper("hash")] Hash hash,
            [Perper("type")] Type type,
            CancellationToken cancellationToken)
        {
            Console.WriteLine("Read: {0}", hash);
            byte[]? serialized;

            // Simulate IPFS's behavior where trying to get a nonexistent object blocks until the object is available.
            while (!HashRegistryStream._storedValues.TryGetValue(hash, out serialized))
            {
                cancellationToken.ThrowIfCancellationRequested();
                await Task.Delay(50, cancellationToken);
            }

            var value = JsonSerializer.Deserialize<Block>(serialized!, ApocryphSerializationOptions.JsonSerializerOptions);
            return value;
        }
    }
}