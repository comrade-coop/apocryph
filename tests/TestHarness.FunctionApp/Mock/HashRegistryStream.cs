using System.Collections.Generic;
using System.Collections.Concurrent;
using System.Security.Cryptography;
using System.Threading.Tasks;
using System.Text.Json;
using Apocryph.Core.Consensus.Blocks;
using Apocryph.Core.Consensus.Serialization;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace TestHarness.FunctionApp.Mock
{
    public class HashRegistryStream
    {
        internal static readonly ConcurrentDictionary<Hash, byte[]> _storedValues = new ConcurrentDictionary<Hash, byte[]>();

        [FunctionName(nameof(HashRegistryStream))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("input")] IAsyncEnumerable<Block> input)
        {
            using var sha256Hash = SHA256.Create();
            await foreach (var value in input)
            {
                var serialized = JsonSerializer.SerializeToUtf8Bytes(value, ApocryphSerializationOptions.JsonSerializerOptions);
                var hash = new Hash(sha256Hash.ComputeHash(serialized));

                _storedValues.TryAdd(hash, serialized);

                // Console.WriteLine("Store: {0}", hash);
            }
        }
    }
}