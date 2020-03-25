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
    public static class CommitCounter
    {
        private class State
        {
            public Dictionary<Hash, HashSet<ValidatorKey>> Commits { get; set; } = new Dictionary<Hash, HashSet<ValidatorKey>>();
            public ValidatorSet ValidatorSet { get; set; }
        }

        [FunctionName(nameof(CommitCounter))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("validatorSetsStream")] IAsyncEnumerable<IHashed<ValidatorSet>> validatorSetsStream,
            [PerperStream("commitsStream")] IAsyncEnumerable<ISigned<Commit>> commitsStream,
            [PerperStream("outputStream")] IAsyncCollector<Hash> outputStream,
            ILogger logger)
        {
            var state = await context.FetchStateAsync<State>() ?? new State();

            await Task.WhenAll(
                validatorSetsStream.ForEachAsync(async validatorSet =>
                {
                    try
                    {
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

                        var wasMoreThanTwoThirds = state.ValidatorSet.IsMoreThanTwoThirds(state.Commits[commit.Value.For]);

                        state.Commits[commit.Value.For].Add(commit.Signer);
                        await context.UpdateStateAsync(state);

                        if (!wasMoreThanTwoThirds &&
                            state.ValidatorSet.IsMoreThanTwoThirds(state.Commits[commit.Value.For]))
                        {
                            await outputStream.AddAsync(commit.Value.For);
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