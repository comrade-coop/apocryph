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
        }

        [FunctionName(nameof(CurrentProposer))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("validatorSet")] ValidatorSet validatorSet,
            [PerperStream("commitsStream")] IAsyncEnumerable<ISigned<Commit>> commitsStream,
            [PerperStream("outputStream")] IAsyncCollector<ValidatorKey> outputStream,
            ILogger logger)
        {
            var state = await context.FetchStateAsync<State>() ?? new State();

            await commitsStream.ForEachAsync(async commit =>
            {
                try
                {
                    if (!state.Commits.ContainsKey(commit.Value.For))
                    {
                        state.Commits[commit.Value.For] = new HashSet<ValidatorKey>();
                    }
                    // TODO: Timeout proposers, rotate proposer only on his own blocks
                    state.Commits[commit.Value.For].Add(commit.Signer);
                    await context.UpdateStateAsync(state);

                    var committed = state.Commits[commit.Value.For]
                        .Select(signer => validatorSet.Weights[signer]).Sum();
                    if (3 * committed > 2 * validatorSet.Total)
                    {
                        validatorSet.AccumulateWeights();
                        var proposer = validatorSet.PopMaxAccumulatedWeight();

                        await outputStream.AddAsync(proposer);
                    }
                }
                catch (Exception e)
                {
                    logger.LogError(e.ToString());
                }
            }, CancellationToken.None);
        }
    }
}