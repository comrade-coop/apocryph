using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Model;
using Microsoft.Azure.WebJobs;
using Microsoft.Extensions.Logging;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.FunctionApp
{
    public static class Proposer
    {
        [FunctionName(nameof(Proposer))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("proposerCommitInjectorStream")] IAsyncEnumerable<IHashed<IAgentStep>> proposerCommitInjectorStream,
            [PerperStream("outputStream")] IAsyncCollector<Proposal> outputStream,
            ILogger logger)
        {
            await proposerCommitInjectorStream.ForEachAsync(async step =>
            {
                try
                {
                    await outputStream.AddAsync(new Proposal{ For = step.Hash });
                }
                catch (Exception e)
                {
                    logger.LogError(e.ToString());
                }
            }, CancellationToken.None);
        }
    }
}