using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Model;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Bindings;
using Perper.WebJobs.Extensions.Model;
using Perper.WebJobs.Extensions.Triggers;

namespace Apocryph.FunctionApp
{
    public static class Committer
    {
        private class State
        {
            public Dictionary<IAgentStep, Dictionary<ValidatorKey, ValidatorSignature>> Commits { get; set; }
        }

        [FunctionName("Committer")]
        public static async Task Run([PerperStreamTrigger] IPerperStreamContext context,
            [PerperStream("self")] ValidatorKey self,
            [PerperStream("validatorSet")] ValidatorSet validatorSet,
            [PerperStream("commitsStream")] IAsyncEnumerable<Commit> commitsStream,
            [PerperStream] IAsyncCollector<(IAgentStep, bool)> outputStream)
        {
            var state = await context.GetState<State>("state");

            await commitsStream.Listen(
                async commit =>
                {
                    state.Commits[commit.For].Add(commit.Signer, commit.Signature);
                    await context.SetState("state", state);

                    var committed = state.Commits[commit.For].Keys
                        .Select(signer => validatorSet.Weights[signer]).Sum();
                    if (3 * committed > 2 * validatorSet.Total)
                    {
                        validatorSet.AccumulateWeights();
                        var proposer = validatorSet.PopMaxAccumulatedWeight();
                        var isProposer = proposer.Equals(self);

                        await outputStream.AddAsync((commit.For, isProposer));
                    }
                },
                CancellationToken.None);
        }
    }
}