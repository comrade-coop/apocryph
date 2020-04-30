// NOTE: File is ignored by .csproj file

using System;
using System.Collections.Generic;
using System.Linq;
using System.Security.Cryptography;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Agent;
using Apocryph.Runtime.FunctionApp.Ipfs;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Runtime.FunctionApp.ValidatorSelection
{
    public static class KingOfTheHillKeySearch
    {
        private class State : KingOfTheHillBase.State
        {
        }

        [FunctionName(nameof(KingOfTheHillKeySearch))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("seenKeysStream")] IAsyncEnumerable<ValidatorKey> seenKeysStream,
            [PerperStream("saltsStream")] IAsyncEnumerable<(int, byte[])> saltsStream,
            [PerperStream("outputStream")] IAsyncCollector<ECParameters> outputStream,
            CancellationToken cancellationToken)
        {
            var state = await context.FetchStateAsync<State>() ?? new State();

            Task.Run(() => // TODO: Run multiple threads in parallel
            {
                while (true)
                {
                    using var dsa = ECDsa.Create(ECCurve.NamedCurves.nistP521);
                    var publicKey = new ValidatorKey{Key = dsa.ExportParameters(false)};
                    if (state.AddKey(publicKey))
                    {
                        var privateKey = dsa.ExportParameters(true);
                        outputStream.AddAsync(privateKey);
                    }
                }
            });

            await KingOfTheHillBase.Run(context, state, seenKeysStream, saltsStream);
        }
    }
}