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
    public static class Committer
    {
        private class State
        {
            public Dictionary<Hash, Dictionary<ValidatorKey, ValidatorSignature>> Commits { get; set; }
        }

        [FunctionName("Committer")]
        public static async Task Run([Perper(Stream = "Committer")] IPerperStreamContext context,
            [Perper("self")] ValidatorKey self,
            [Perper("validatorSet")] ValidatorSet validatorSet,
            [Perper("commitsStream")] IAsyncEnumerable<Signed<Commit>> commitsStream,
            [Perper("outputStream")] IAsyncCollector<(Hash, bool)> outputStream)
        {
            var state = context.GetState<State>("state");

            await commitsStream.Listen(
                async commit =>
                {
                    state.Commits[commit.Value.For].Add(commit.Signer, commit.Signature);
                    await context.SaveState("state", state);

                    var committed = state.Commits[commit.Value.For].Keys
                        .Select(signer => validatorSet.Weights[signer]).Sum();
                    if (3 * committed > 2 * validatorSet.Total)
                    {
                        validatorSet.AccumulateWeights();
                        var proposer = validatorSet.PopMaxAccumulatedWeight();
                        var isProposer = proposer.Equals(self);

                        await outputStream.AddAsync((commit.Value.For, isProposer));
                    }
                },
                CancellationToken.None);
        }
    }
}