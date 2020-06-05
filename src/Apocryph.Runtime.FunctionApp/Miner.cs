using System.Collections.Generic;
using System.Security.Cryptography;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Core.Consensus.VirtualNodes;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Runtime.FunctionApp
{
    public static class Miner
    {
        [FunctionName(nameof(Miner))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("chains")] IDictionary<byte[], string> chains,
            [PerperStream("output")] IAsyncCollector<(PrivateKey, string)> output,
            CancellationToken cancellationToken)
        {

            using var dsa = ECDsa.Create();
            // TODO: Maybe run multiple threads in parallel
            while (!cancellationToken.IsCancellationRequested)
            {
                dsa.GenerateKey(PrivateKey.Curve); // Pass chains.keys as prefixes
                var privateKey = new PrivateKey(dsa.ExportParameters(true));
                await output.AddAsync((privateKey, chains[GetPrefix(privateKey)]));
            }
        }

        private static byte[] GetPrefix(PrivateKey privateKey)
        {
            return default!; //prefix of the private key
        }
    }
}