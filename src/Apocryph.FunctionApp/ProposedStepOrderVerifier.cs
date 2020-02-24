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
    public static class ProposedStepOrderVerifier
    {
        private class State
        {
            public Hash CurrentStep { get; set; } = new Hash {Bytes = new byte[] {}};
            public List<Hash> ValidatorSets { get; set; } = new List<Hash> {new Hash { Bytes = new byte[]{} }};
        }

        [FunctionName(nameof(ProposedStepOrderVerifier))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("stepsStream")] IAsyncEnumerable<IHashed<IAgentStep>> stepsStream,
            [PerperStream("proposedStepsStream")] IAsyncEnumerable<IHashed<IAgentStep>> proposedStepsStream,
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
                        var i = 0;
                        for (; i < state.ValidatorSets.Count; i++)
                        {
                            if (step.Value.PreviousValidatorSet == state.ValidatorSets[i])
                            {
                                break;
                            }
                        }
                        state.ValidatorSets.RemoveRange(0, i);
                        state.CurrentStep = step.Hash;
                        await context.UpdateStateAsync(state);
                    }
                    catch (Exception e)
                    {
                        logger.LogError(e.ToString());
                    }
                }, CancellationToken.None),

                proposedStepsStream.ForEachAsync(async proposedStep =>
                {
                    try
                    {
                        if (proposedStep.Value.Previous != state.CurrentStep)
                        {
                            return;
                        }

                        for (var i = 0; i < state.ValidatorSets.Count; i++)
                        {
                            if (proposedStep.Value.PreviousValidatorSet == state.ValidatorSets[i])
                            {
                                await outputStream.AddAsync(proposedStep);
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