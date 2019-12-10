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
    public static class Consensus
    {
        private class State
        {
            public Dictionary<IAgentStep, Dictionary<ValidatorKey, ValidatorSignature>> Votes { get; set; }
        }

        [FunctionName("Consensus")]
        public static async Task Run([PerperStreamTrigger] IPerperStreamContext context,
            [PerperStream("validatorSet")] ValidatorSet validatorSet,
            [PerperStream("votesStream")] IAsyncEnumerable<Vote> votesStream,
            [PerperStream] IAsyncCollector<Commit> outputStream)
        {
            var state = await context.GetState<State>("state");

            await votesStream.Listen(
                async vote =>
                {
                    state.Votes[vote.For].Add(vote.Signer, vote.Signature);
                    await context.SetState("state", state);

                    var voted = state.Votes[vote.For].Keys
                        .Select(signer => validatorSet.Weights[signer]).Sum();
                    if (3 * voted > 2 * validatorSet.Total)
                    {
                        await outputStream.AddAsync(new Commit {For = vote.For});
                    }
                },
                CancellationToken.None);
        }
    }
}