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
    public static class StepVerifiedStepGetter
    {
        [FunctionName(nameof(StepVerifiedStepGetter))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("stepSignatureVerifierStream")] IAsyncEnumerable<IHashed<IAgentStep>> stepSignatureVerifierStream,
            [PerperStream("outputStream")] IAsyncCollector<Hash> outputStream)
        {
            await stepSignatureVerifierStream.ForEachAsync(async step =>
            {
                await outputStream.AddAsync(step.Value.Previous);
            }, CancellationToken.None);
        }
    }
}