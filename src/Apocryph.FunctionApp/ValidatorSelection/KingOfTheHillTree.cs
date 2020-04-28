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
using Microsoft.Extensions.Logging;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.FunctionApp.ValidatorSelection
{
    public static class KingOfTheHillTree
    {
        private class State
        {
            // public ValidatorTree Tree { get; set; } = new ValidatorTree();
        }

        [FunctionName(nameof(KingOfTheHillTree))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("seenKeysStream")] IAsyncEnumerable<ValidatorKey> seenKeysStream,
            [PerperStream("agentDesciptorsStream")] IAsyncEnumerable<AgentDescriptor> agentDesciptorsStream,
            [PerperStream("outputStream")] IAsyncCollector<List<ValidatorKey>> outputStream,
            ILogger logger)
        {
            var state = await context.FetchStateAsync<State>() ?? new State();

            await Task.WhenAll(
                seenKeysStream.ForEachAsync(async key =>
                {
                    try
                    {
                        // state.Tree.AddValidator(key);
                        await context.UpdateStateAsync(state);
                    }
                    catch (Exception e)
                    {
                        logger.LogError(e.ToString());
                    }
                }, CancellationToken.None),

                agentDesciptorsStream.ForEachAsync(async agentDescriptor =>
                {
                    try
                    {
                        // await outputStream.AddAsync(state.Tree.GetValidators(agentDescriptor.ValidatorPrefix, agentDescriptor.ValidatorDepth).ToList());
                    }
                    catch (Exception e)
                    {
                        logger.LogError(e.ToString());
                    }
                }, CancellationToken.None));
        }
    }
}