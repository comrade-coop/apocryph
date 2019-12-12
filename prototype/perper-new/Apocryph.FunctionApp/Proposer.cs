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
    public static class Proposer
    {
        private class State
        {
            public Dictionary<Hash, Dictionary<ValidatorKey, ValidatorSignature>> Commits { get; set; }
        }

        [FunctionName("Proposer")]
        public static async Task Run([PerperTrigger("Proposer")] IPerperStreamContext context,
            [Perper("commitsStream")] IAsyncEnumerable<Signed<Commit>> commitsStream,
            [Perper("proposerRuntimeStream")] IAsyncEnumerable<IAgentStep> proposerRuntimeStream,
            [Perper("outputStream")] IAsyncCollector<IAgentStep> outputStream)
        {
            var state = context.GetState<State>();

            await Task.WhenAll(
                commitsStream.ForEachAsync(async commit =>
                {
                    state.Commits[commit.Value.For].Add(commit.Signer, commit.Signature);
                    await context.SaveState();
                }, CancellationToken.None),

                proposerRuntimeStream.ForEachAsync(async step =>
                {
                    // FIXME: Should probably block until there are enough signatures (as we cannot be sure that the other stream would collect them in time)
                    step.CommitSignatures = state.Commits[step.Previous];

                    await outputStream.AddAsync(step);
                }, CancellationToken.None));
        }
    }
}