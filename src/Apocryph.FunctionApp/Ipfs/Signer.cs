using System;
using System.Collections.Generic;
using System.Security.Cryptography;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Model;
using Apocryph.FunctionApp.Ipfs;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.FunctionApp
{
    public static class Signer
    {
        [FunctionName(nameof(Signer))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("self")] ValidatorKey self,
            [Perper("privateKey")] ECParameters privateKey,
            [PerperStream("dataStream")] IAsyncEnumerable<object> dataStream,
            [PerperStream("outputStream")] IAsyncCollector<ISigned<object>> outputStream)
        {

            await dataStream.ForEachAsync(async item =>
                {
                    var bytes = IpfsJsonSettings.ObjectToBytes(item);
                    var signature = ValidatorKey.GenerateSignature(privateKey, bytes);

                    await outputStream.AddAsync(Signed.Create(item, self, signature));
                }, CancellationToken.None);
        }
    }
}