// NOTE: File is ignored by .csproj file

using System;
using System.Collections.Generic;
using System.Linq;
using System.Security.Cryptography;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Agent;
using Apocryph.FunctionApp.Command;
using Apocryph.FunctionApp.Model;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.FunctionApp.ValidatorSelection
{
    public static class KingOfTheHillKeySearch
    {
        [FunctionName(nameof(KingOfTheHillKeySearch))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("seenKeysStream")] IAsyncEnumerable<ValidatorKey> seenKeysStream,
            [PerperStream("outputStream")] IAsyncCollector<ECParameters> outputStream,
            CancellationToken cancellationToken)
        {
        }
    }
}