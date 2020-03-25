using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Model;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.FunctionApp
{
    public static class ValidatorRuntimeInput
    {
        [FunctionName(nameof(ValidatorRuntimeInput))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("validatorFilterStream")] IAsyncEnumerable<IHashed<AgentOutput>> validatorFilterStream,
            [PerperStream("outputStream")] IAsyncCollector<Hash> outputStream)
        {
            await validatorFilterStream.ForEachAsync(async output =>
            {
                await outputStream.AddAsync(output.Value.Previous);
            }, CancellationToken.None);
        }
    }
}