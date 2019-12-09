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
    public static class Validator
    {
        private class State
        {
            public string Proposer { get; set; }
            public Dictionary<(AgentInput, AgentOutput), HashSet<string>> Commits { get; set; }
        }

        [FunctionName("Validator")]
        public static async Task Run([PerperStreamTrigger] IPerperStreamContext context,
            [PerperStream("self")] string self,
            [PerperStream("validatorSet")] ValidatorSet validatorSet,
            [PerperStream("commitsStream")] IPerperStream<Commit> commitsStream,
            [PerperStream("proposalsStream")] IPerperStream<(AgentInput, AgentOutput)> proposalsStream,
            [PerperStream] IAsyncCollector<object> outputStream)
        {
            var state = await context.GetState<State>("state");

            await Task.WhenAll(
                commitsStream.Listen(async commit =>
                {
                    state.Commits[(commit.Input, commit.Output)].Add(commit.Signer);

                    var voted = state.Commits[(commit.Input, commit.Output)]
                        .Select(signer => validatorSet.Weights[signer]).Sum();

                    if (3 * voted > 2 * validatorSet.Total)
                    {
                        validatorSet.AccumulateWeights();
                        state.Proposer = validatorSet.PopMaxAccumulatedWeight();
                    }

                    await context.SetState("state", state);
                }, CancellationToken.None),

                proposalsStream.Listen(async proposal =>
                {
                    if (state.Proposer != self) // state.Proposer == proposal.Signer
                    {
                        await outputStream.AddAsync(proposal.Item1);
                    }
                }, CancellationToken.None));
        }
    }
}