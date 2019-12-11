using System.Threading;
using System.Threading.Tasks;
using System.Collections.Generic;
using Apocryph.FunctionApp.Model;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.FunctionApp
{
    public static class Proposer
    {
        private class State
        {
            public Dictionary<Hash, Dictionary<ValidatorKey, ValidatorSignature>> Commits { get; set; }
        }

        [FunctionName("Proposer")]
        public static async Task Run([Perper(Stream = "Proposer")] IPerperStreamContext context,
            [Perper("commitsStream")] IAsyncEnumerable<Commit> commitsStream,
            [Perper("runtimeStream")] IAsyncEnumerable<(IAgentStep, bool)> runtimeStream,
            [Perper("outputStream")] IAsyncCollector<IAgentStep> outputStream)
        {
            var state = context.GetState<State>("state");

            await Task.WhenAll(
                commitsStream.Listen(async commit =>
                {
                    state.Commits[commit.ForHash].Add(commit.Signer, commit.Signature);
                    await context.SaveState("state", state);
                }, CancellationToken.None),

                runtimeStream.Listen(async item =>
                {
                    var (step, isProposal) = item;
                    if (isProposal)
                    {
                        // FIXME: Should probably wait to accumulate enough signatures (as we cannot be sure that the other stream would collect them in time)
                        step.CommitSignatures = state.Commits[step.PreviousHash];

                        await outputStream.AddAsync(step);
                    }
                }, CancellationToken.None));
        }
    }
}