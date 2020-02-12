using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Agent;
using Apocryph.FunctionApp.Command;
using Apocryph.FunctionApp.Ipfs;
using Apocryph.FunctionApp.Model;
using Microsoft.Azure.WebJobs;
using Microsoft.Extensions.Logging;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.FunctionApp
{
    public static class StepValidatorSetSplitter
    {
        [FunctionName(nameof(StepValidatorSetSplitter))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("stepsStream")] IAsyncEnumerable<IHashed<IAgentStep>> stepsStream,
            [PerperStream("outputStream")] IAsyncCollector<Hash> outputStream,
            ILogger logger)
        {
            await stepsStream.ForEachAsync(async step =>
            {
                await outputStream.AddAsync(step.Value.PreviousValidatorSet);
            }, CancellationToken.None);
        }
    }
}