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
    public static class Validator
    {
        private class State
        {
            public ValidatorKey Proposer { get; set; }
            public IAgentStep CurrentStep { get; set; }
            public Dictionary<IAgentStep, HashSet<ValidatorKey>> Commits { get; set; }
        }

        [FunctionName("Validator")]
        public static async Task Run([Perper(Stream = "Validator")] IPerperStreamContext context,
            [Perper("validatorSet")] ValidatorSet validatorSet,
            [Perper("commitsStream")] IAsyncEnumerable<Commit> commitsStream,
            [Perper("proposalsStream")] IAsyncEnumerable<IAgentStep> proposalsStream,
            [Perper("outputStream")] IAsyncCollector<IAgentStep> outputStream)
        {
            var state = context.GetState<State>("state");

            await Task.WhenAll(
                commitsStream.Listen(async commit =>
                {
                    state.Commits[commit.For].Add(commit.Signer);

                    var committed = state.Commits[commit.For]
                        .Select(signer => validatorSet.Weights[signer]).Sum();

                    if (3 * committed > 2 * validatorSet.Total)
                    {
                        validatorSet.AccumulateWeights();
                        state.Proposer = validatorSet.PopMaxAccumulatedWeight();
                        state.CurrentStep = commit.For; // TODO: Commit in order
                    }

                    await context.SaveState("state", state);
                }, CancellationToken.None),

                proposalsStream.Listen(async proposal =>
                {
                    if (state.Proposer.Equals(proposal.Signer) && state.CurrentStep == proposal.Previous)
                    {
                        await outputStream.AddAsync(proposal.Previous);
                    }
                }, CancellationToken.None));
        }
    }
}