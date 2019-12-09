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
            public Dictionary<(AgentInput, AgentOutput), HashSet<string>> Commits { get; set; }
        }

        [FunctionName("Committer")]
        public static async Task Run([PerperStreamTrigger] IPerperStreamContext context,
            [PerperStream("self")] string self, // Or should it be "ownPublicKey" / "selfSigner"?
            [PerperStream("validatorSet")] ValidatorSet validatorSet,
            [PerperStream("commitsStream")] IPerperStream<Commit> commitsStream,
            [PerperStream] IAsyncCollector<(AgentOutput, bool)> outputStream)
        {
            var state = await context.GetState<State>("state");

            await commitsStream.Listen(
                async commit =>
                {
                    state.Commits[(commit.Input, commit.Output)].Add(commit.Signer);
                    await context.SetState("state", state);

                    var voted = state.Commits[(commit.Input, commit.Output)]
                        .Select(signer => validatorSet.Weights[signer]).Sum();
                    if (3 * voted > 2 * validatorSet.Total)
                    {
                        validatorSet.AccumulateWeights();
                        var proposer = validatorSet.PopMaxAccumulatedWeight();
                        var isProposer = (proposer == self);

                        await outputStream.AddAsync((commit.Output, isProposer));
                    }
                },
                CancellationToken.None);
        }
    }
}