// NOTE: File is ignored by .csproj file

using System.Security.Cryptography;
using System.Threading;
using System.Threading.Tasks;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Runtime.FunctionApp.ValidatorSelection
{
    public static class Miner
    {
        [FunctionName(nameof(Miner))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("output")] IAsyncCollector<PrivateKey> output,
            CancellationToken cancellationToken)
        {

            using var dsa = ECDsa.Create();
            // TODO: Maybe run multiple threads in parallel
            while (!cancellationToken.IsCancellationRequested)
            {
                dsa.GenerateKey(PrivateKey.Curve);
                var privateKey = new PrivateKey(dsa.ExportParameters(true));
                await output.AddAsync(privateKey);
            }
        }
    }
}