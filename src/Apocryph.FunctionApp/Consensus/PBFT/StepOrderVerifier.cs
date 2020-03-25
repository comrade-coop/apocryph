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
    public static class StepOrderVerifier
    {
        private class State
        {
            public Hash CurrentStep { get; set; } = new Hash {Bytes = new byte[] {}};
            public List<Hash> ValidatorSets { get; set; } = new List<Hash> {new Hash { Bytes = new byte[]{} }};
        }

        [FunctionName(nameof(StepOrderVerifier))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("stepsStream")] IAsyncEnumerable<IHashed<IAgentStep>> stepsStream,
            [PerperStream("validatorSetsStream")] IAsyncEnumerable<IHashed<ValidatorSet>> validatorSetsStream,
            [PerperStream("outputStream")] IAsyncCollector<IHashed<IAgentStep>> outputStream,
            ILogger logger)
        {
            var state = await context.FetchStateAsync<State>() ?? new State();

            await Task.WhenAll(
                validatorSetsStream.ForEachAsync(async validatorSet =>
                {
                    state.ValidatorSets.Add(validatorSet.Hash);

                    await context.UpdateStateAsync(state);
                }, CancellationToken.None),

                stepsStream.ForEachAsync(async step =>
                {
                    try
                    {
                        if (step.Value.Previous != state.CurrentStep)
                        {
                            return;
                        }

                        // Find first matching validator set, drop all before it
                        for (var i = 0; i < state.ValidatorSets.Count; i++)
                        {
                            if (step.Value.PreviousValidatorSet == state.ValidatorSets[i])
                            {
                                state.ValidatorSets.RemoveRange(0, i);
                                state.CurrentStep = step.Hash;
                                await outputStream.AddAsync(step);
                                await context.UpdateStateAsync(state);
                                break;
                            }
                        }
                    }
                    catch (Exception e)
                    {
                        logger.LogError(e.ToString());
                    }
                }, CancellationToken.None));
        }
    }
}