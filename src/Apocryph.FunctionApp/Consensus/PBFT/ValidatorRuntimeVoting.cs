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
    public static class ValidatorRuntimeVoting
    {
        private class State
        {
            public Dictionary<Hash, Hash> ExpectedNextSteps { get; } = new Dictionary<Hash, Hash>();
        }

        [FunctionName(nameof(ValidatorRuntimeVoting))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("validatorRuntimeOutputStream")] IAsyncEnumerable<IHashed<AgentOutput>> validatorRuntimeOutputStream,
            [PerperStream("validatorFilterStream")] IAsyncEnumerable<IHashed<AgentOutput>> validatorFilterStream,
            [PerperStream("genesisStream")] IAsyncEnumerable<IHashed<AgentOutput>> genesisStream,
            [PerperStream("outputStream")] IAsyncCollector<Vote> outputStream,
            ILogger logger)
        {
            var state = await context.FetchStateAsync<State>() ?? new State();

            await Task.WhenAll(
                validatorFilterStream.ForEachAsync(async proposal =>
                {
                    try
                    {
                        state.ExpectedNextSteps[proposal.Value.Previous] = proposal.Hash;

                        await context.UpdateStateAsync(state);
                    }
                    catch (Exception e)
                    {
                        logger.LogError(e.ToString());
                    }
                }, CancellationToken.None),

                validatorRuntimeOutputStream.ForEachAsync(async step =>
                {
                    try
                    {
                        if (state.ExpectedNextSteps[step.Value.Previous] == step.Hash)
                        {
                            await outputStream.AddAsync(new Vote { For = step.Hash });
                        }
                        else
                        {
                            logger.LogWarning("Got {hash}, expected {expected}", step.Hash, state.ExpectedNextSteps[step.Value.Previous]);
                        }
                    }
                    catch (Exception e)
                    {
                        logger.LogError(e.ToString());
                    }
                }, CancellationToken.None),

                genesisStream.ForEachAsync(async step =>
                {
                    try
                    {
                        await outputStream.AddAsync(new Vote { For = step.Hash });
                    }
                    catch (Exception e)
                    {
                        logger.LogError(e.ToString());
                    }
                }, CancellationToken.None));
        }
    }
}