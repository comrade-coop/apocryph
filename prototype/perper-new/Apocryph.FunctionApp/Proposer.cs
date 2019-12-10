using System.Threading;
using System.Threading.Tasks;
using System.Collections.Generic;
using Apocryph.FunctionApp.Model;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Bindings;
using Perper.WebJobs.Extensions.Model;
using Perper.WebJobs.Extensions.Triggers;

namespace Apocryph.FunctionApp
{
    public static class Proposer
    {
        private class State
        {
            public Dictionary<IAgentStep, Dictionary<ValidatorKey, ValidatorSignature>> Commits { get; set; }
        }

        [FunctionName("Proposer")]
        public static async Task Run([PerperStreamTrigger] IPerperStreamContext context,
            [PerperStream("commitsStream")] IAsyncEnumerable<Commit> commitsStream,
            [PerperStream("runtimeStream")] IAsyncEnumerable<(IAgentStep, bool)> runtimeStream,
            [PerperStream] IAsyncCollector<IAgentStep> outputStream)
        {
            var state = await context.GetState<State>("state");

            await Task.WhenAll(
                commitsStream.Listen(async commit =>
                {
                    state.Commits[commit.For].Add(commit.Signer, commit.Signature);
                    await context.SetState("state", state);
                }, CancellationToken.None),

                runtimeStream.Listen(async item =>
                {
                    var (step, isProposal) = item;
                    if (isProposal)
                    {
                        // FIXME: Should probably wait to accumulate enough signatures (as we cannot be sure that the other stream would collect them in time)
                        step.CommitSignatures = state.Commits[step.Previous];

                        await outputStream.AddAsync(step);
                    }
                }, CancellationToken.None));
        }
    }
}