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
    public static class ProposerFilter
    {
        private class State
        {
            public IHashed<IAgentStep>? LastCommit { get; set; } = null;
            public ValidatorKey? LastProposer { get; set; } = null;
        }

        [FunctionName(nameof(ProposerFilter))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("self")] ValidatorKey self,
            [PerperStream("syncStream")] IAsyncEnumerable<IHashed<IAgentStep>> syncStream,
            [PerperStream("currentProposerStream")] IAsyncEnumerable<ValidatorKey> currentProposerStream,
            [PerperStream("outputStream")] IAsyncCollector<IHashed<IAgentStep>> outputStream)
        {
            var state = await context.FetchStateAsync<State>() ?? new State();

            async Task processPending()
            {
                if (state.LastProposer != null && state.LastProposer.Equals(self) && state.LastCommit != null)
                {
                    await outputStream.AddAsync(state.LastCommit);
                    state.LastCommit = null; // Do not produce the same commit twice
                    state.LastProposer = null; // Do not produce the next commit before the proposer is updated
                }
            };

            await Task.WhenAll(
                currentProposerStream.ForEachAsync(async currentProposer =>
                {
                    state.LastProposer = currentProposer;
                    await processPending();
                    await context.UpdateStateAsync(state);
                }, CancellationToken.None),

                syncStream.ForEachAsync(async commit =>
                {
                    state.LastCommit = commit;
                    await processPending();
                    await context.UpdateStateAsync(state);
                }, CancellationToken.None));
        }
    }
}