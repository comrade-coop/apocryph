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
    public static class Consensus
    {
        private class State
        {
            public Dictionary<Hash, HashSet<ValidatorKey>> Votes { get; set; } = new Dictionary<Hash, HashSet<ValidatorKey>>();
            public ValidatorSet ValidatorSet { get; set; }
        }

        [FunctionName(nameof(Consensus))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("validatorSetsStream")] IAsyncEnumerable<IHashed<ValidatorSet>> validatorSetsStream,
            [PerperStream("votesStream")] IAsyncEnumerable<ISigned<Vote>> votesStream,
            [PerperStream("outputStream")] IAsyncCollector<Commit> outputStream,
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

                votesStream.ForEachAsync(async vote =>
                {
                    if (!state.Votes.ContainsKey(vote.Value.For))
                    {
                        state.Votes[vote.Value.For] = new HashSet<ValidatorKey>();
                    }
                    var wasMoreThanTwoThirds = state.ValidatorSet.IsMoreThanTwoThirds(state.Votes[vote.Value.For]);

                    state.Votes[vote.Value.For].Add(vote.Signer);
                    await context.UpdateStateAsync(state);

                    if (!wasMoreThanTwoThirds &&
                        state.ValidatorSet.IsMoreThanTwoThirds(state.Votes[vote.Value.For]))
                    {
                        await outputStream.AddAsync(new Commit {For = vote.Value.For});
                    }
                }, CancellationToken.None));
        }
    }
}