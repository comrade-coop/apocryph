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
            [PerperStream("validatorSetsStream")] IAsyncEnumerable<IHashed<ValidatorSet>> validatorSetsStream,
            [PerperStream("commitsStream")] IAsyncEnumerable<ISigned<Commit>> commitsStream,
            [PerperStream("outputStream")] IAsyncCollector<ValidatorKey> outputStream,
            ILogger logger)
        {
            var state = await context.FetchStateAsync<State>() ?? new State();

            await Task.WhenAll(
                validatorSetsStream.ForEachAsync(async validatorSet =>
                {
                    try
                    {
                        if (state.ValidatorSet == null) // HACK
                        {
                            // state.ValidatorSet.AccumulateWeights();
                            // state.ValidatorSet.PopMaxAccumulatedWeight();
                            var initialProposer = validatorSet.Value.Weights.Keys.First();
                            await outputStream.AddAsync(initialProposer);
                        }
                        state.ValidatorSet = validatorSet.Value;
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
                        var wasMoreThanTwoThirds = state.ValidatorSet.IsMoreThanTwoThirds(state.Commits[commit.Value.For]);

                        state.Commits[commit.Value.For].Add(commit.Signer);
                        if (!wasMoreThanTwoThirds &&
                            state.ValidatorSet.IsMoreThanTwoThirds(state.Commits[commit.Value.For]))
                        {
                            // state.ValidatorSet.AccumulateWeights();
                            // state.ValidatorSet.PopMaxAccumulatedWeight();
                            var proposer = state.ValidatorSet.Weights.Keys.First();

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