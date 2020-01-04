using System;
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
        [FunctionName(nameof(Signer))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("self")] ValidatorKey self,
            [Perper("privateKey")] ECParameters privateKey,
            [PerperStream("dataStream")] IAsyncEnumerable<IHashed<object>> dataStream,
            [PerperStream("outputStream")] IAsyncCollector<ISigned<object>> outputStream)
        {

            await dataStream.ForEachAsync(async hashed =>
                {
                    using var ecdsa = ECDsa.Create(privateKey);
                    var signature = new ValidatorSignature
                    {
                        Bytes = ecdsa.SignHash(hashed.Hash.Bytes)
                    };

                    var signedType = typeof(Signed<>).MakeGenericType(hashed.GetType().GenericTypeArguments[0]);
                    var signed = (ISigned<object>)Activator.CreateInstance(signedType, hashed, self, signature);

                    await outputStream.AddAsync(signed);
                }, CancellationToken.None);
        }
    }
}