using System.Collections.Generic;
using System.Security.Cryptography;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Model;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.FunctionApp
{
    public static class Signer
    {
        [FunctionName("Signer")]
        public static async Task Run([PerperTrigger("Signer")] IPerperStreamContext context,
            [Perper("self")] ValidatorKey self,
            [Perper("privateKey")] ECParameters privateKey,
            [Perper("dataStream")] IAsyncEnumerable<Hashed<object>> dataStream,
            [Perper("outputStream")] IAsyncCollector<Signed<object>> outputStream)
        {

            await dataStream.ForEachAsync(async item =>
                {
                    byte[] signature;
                    using (var ecdsa = ECDsa.Create(privateKey))
                    {
                        signature = ecdsa.SignHash(item.Hash.Bytes);
                    }

                    await outputStream.AddAsync(new Signed<object>(item, self, new ValidatorSignature {Bytes = signature}));
                }, CancellationToken.None);
        }
    }
}