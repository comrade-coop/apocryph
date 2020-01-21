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
    public static class CurrentProposer
    {
        private class State
        {
            public Dictionary<Hash, HashSet<ValidatorKey>> Commits { get; set; } = new Dictionary<Hash, HashSet<ValidatorKey>>();
            public ValidatorSet? ValidatorSet { get; set; } = null;
        }

        [FunctionName(nameof(CurrentProposer))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("validatorSetStream")] IAsyncEnumerable<ValidatorSet> validatorSetStream,
            [PerperStream("commitsStream")] IAsyncEnumerable<ISigned<Commit>> commitsStream,
            [PerperStream("outputStream")] IAsyncCollector<ValidatorKey> outputStream,
            ILogger logger)
        {
            var state = await context.FetchStateAsync<State>() ?? new State();

            await Task.WhenAll(
                validatorSetStream.ForEachAsync(async validatorSet =>
                {
                    try
                    {
                        if (state.ValidatorSet == null) // HACK
                        {
                            state.ValidatorSet.AccumulateWeights();
                            var initialProposer = state.ValidatorSet.PopMaxAccumulatedWeight();
                            await outputStream.AddAsync(initialProposer);
                        }
                        state.ValidatorSet = validatorSet;
                        await context.UpdateStateAsync(state);
                    }
                    catch (Exception e)
                    {
                        logger.LogError(e.ToString());
                    }
                }, CancellationToken.None),

                commitsStream.ForEachAsync(async commit =>
                {
                    try
                    {
                        if (!state.Commits.ContainsKey(commit.Value.For))
                        {
                            state.Commits[commit.Value.For] = new HashSet<ValidatorKey>();
                        }
                        // TODO: Timeout proposers, rotate proposer only on his own blocks
                        state.Commits[commit.Value.For].Add(commit.Signer);

                        var committed = state.Commits[commit.Value.For]
                            .Select(signer => state.ValidatorSet.Weights[signer]).Sum();
                        if (3 * committed > 2 * state.ValidatorSet.Total && 3 * committed - state.ValidatorSet.Weights[commit.Signer] <= 2 * state.ValidatorSet.Total)
                        {
                            state.ValidatorSet.AccumulateWeights();
                            var proposer = state.ValidatorSet.PopMaxAccumulatedWeight();

                            await outputStream.AddAsync(proposer);
                        }
                        await context.UpdateStateAsync(state);
                    }
                    catch (Exception e)
                    {
                        logger.LogError(e.ToString());
                    }
                }, CancellationToken.None));
        }
    }
}