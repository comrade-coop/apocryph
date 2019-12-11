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
            public Dictionary<Hash, Dictionary<ValidatorKey, ValidatorSignature>> Votes { get; set; }
        }

        [FunctionName("Consensus")]
        public static async Task Run([Perper(Stream = "Consensus")] IPerperStreamContext context,
            [Perper("validatorSet")] ValidatorSet validatorSet,
            [Perper("votesStream")] IAsyncEnumerable<Vote> votesStream,
            [Perper("outputStream")] IAsyncCollector<Commit> outputStream)
        {
            var state = context.GetState<State>("state");

            await votesStream.Listen(
                async vote =>
                {
                    state.Votes[vote.ForHash].Add(vote.Signer, vote.Signature);
                    await context.SaveState("state", state);

                    var voted = state.Votes[vote.ForHash].Keys
                        .Select(signer => validatorSet.Weights[signer]).Sum();
                    if (3 * voted > 2 * validatorSet.Total)
                    {
                        await outputStream.AddAsync(new Commit {ForHash = vote.ForHash});
                    }
                },
                CancellationToken.None);
        }
    }
}