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
    public static class Voting
    {
        private class State
        {
            public Dictionary<Hash, Hash> ExpectedNextSteps { get; set; }
        }

        [FunctionName(nameof(Voting))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("validatorRuntimeStream")] IAsyncEnumerable<Hashed<IAgentStep>> validatorRuntimeStream,
            [PerperStream("validatorFilterStream")] IAsyncEnumerable<Signed<IAgentStep>> validatorFilterStream,
            [PerperStream("outputStream")] IAsyncCollector<Vote> outputStream)
        {
            var state = await context.FetchStateAsync<State>() ?? new State();

            await Task.WhenAll(
                validatorFilterStream.ForEachAsync(async proposal =>
                {
                    state.ExpectedNextSteps[proposal.Value.Previous] = proposal.Hash;

                    await context.UpdateStateAsync(state);
                }, CancellationToken.None),

                validatorRuntimeStream.ForEachAsync(async step =>
                {
                    if (state.ExpectedNextSteps[step.Value.Previous] == step.Hash)
                    {
                        await outputStream.AddAsync(new Vote { For = step.Hash });
                    }
                }, CancellationToken.None));
        }
    }
}