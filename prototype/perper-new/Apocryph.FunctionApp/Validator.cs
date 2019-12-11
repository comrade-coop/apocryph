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
            public Hash CurrentStep { get; set; }
            public Dictionary<Hash, HashSet<ValidatorKey>> Commits { get; set; }
        }

        [FunctionName("Validator")]
        public static async Task Run([Perper(Stream = "Validator")] IPerperStreamContext context,
            [Perper("validatorSet")] ValidatorSet validatorSet,
            [Perper("commitsStream")] IAsyncEnumerable<Signed<Commit>> commitsStream,
            [Perper("proposalsStream")] IAsyncEnumerable<Signed<IAgentStep>> proposalsStream,
            [Perper("outputStream")] IAsyncCollector<Hash> outputStream)
        {
            var state = context.GetState<State>("state");

            await Task.WhenAll(
                commitsStream.Listen(async commit =>
                {
                    state.Commits[commit.Value.For].Add(commit.Signer);

                    var committed = state.Commits[commit.Value.For]
                        .Select(signer => validatorSet.Weights[signer]).Sum();

                    if (3 * committed > 2 * validatorSet.Total)
                    {
                        validatorSet.AccumulateWeights();
                        state.Proposer = validatorSet.PopMaxAccumulatedWeight();
                        state.CurrentStep = commit.Value.For; // TODO: Commit in order
                    }

                    await context.SaveState("state", state);
                }, CancellationToken.None),

                proposalsStream.Listen(async proposal =>
                {
                    if (state.Proposer.Equals(proposal.Signer) && state.CurrentStep == proposal.Value.Previous)
                    {
                        await outputStream.AddAsync(proposal.Value.Previous);
                    }
                }, CancellationToken.None));
        }
    }
}