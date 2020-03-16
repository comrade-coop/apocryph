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
    public static class AgentZeroValidatorSetsSplitter
    {
        [FunctionName(nameof(AgentZeroValidatorSetsSplitter))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("validatorSetsStream")] IAsyncEnumerable<Dictionary<string, IHashed<ValidatorSet>>> validatorSetsStream,
            [Perper("agentId")] string agentId,
            [PerperStream("outputStream")] IAsyncCollector<IHashed<ValidatorSet>> outputStream,
            ILogger logger)
        {
            await validatorSetsStream.ForEachAsync(async validatorSets =>
            {
                if (validatorSets.TryGetValue(agentId, out var value))
                {
                    await outputStream.AddAsync(value);
                }
            }, CancellationToken.None);
        }
    }
}