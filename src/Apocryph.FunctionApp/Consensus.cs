using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Model;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.FunctionApp
{
    public static class Consensus
    {
        private class State
        {
            public Dictionary<Hash, Dictionary<ValidatorKey, ValidatorSignature>> Votes { get; set; } = new Dictionary<Hash, Dictionary<ValidatorKey, ValidatorSignature>>();
        }

        [FunctionName(nameof(Consensus))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("validatorSet")] ValidatorSet validatorSet,
            [PerperStream("votesStream")] IAsyncEnumerable<ISigned<Vote>> votesStream,
            [PerperStream("outputStream")] IAsyncCollector<Commit> outputStream)
        {
            var state = await context.FetchStateAsync<State>() ?? new State();

            await votesStream.ForEachAsync(async vote =>
            {
                if (!state.Votes.ContainsKey(vote.Value.For))
                {
                    state.Votes[vote.Value.For] = new Dictionary<ValidatorKey, ValidatorSignature>();
                }
                state.Votes[vote.Value.For].Add(vote.Signer, vote.Signature);
                await context.UpdateStateAsync(state);

                var voted = state.Votes[vote.Value.For].Keys
                    .Select(signer => validatorSet.Weights[signer]).Sum();
                if (3 * voted > 2 * validatorSet.Total && 3 * voted - validatorSet.Weights[vote.Signer] <= 2 * validatorSet.Total)
                {
                    await outputStream.AddAsync(new Commit {For = vote.Value.For});
                }
            }, CancellationToken.None);
        }
    }
}