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
            public Signed<IAgentStep> LastCommit { get; set; }
        }

        [FunctionName("ProposerFilter")]
        public static async Task Run([PerperStreamTrigger("ProposerFilter")] IPerperStreamContext context,
            [Perper("self")] ValidatorKey self,
            [Perper("committerStream")] IAsyncEnumerable<Signed<IAgentStep>> committerStream,
            [Perper("currentProposerStream")] IAsyncEnumerable<ValidatorKey> currentProposerStream,
            [Perper("outputStream")] IAsyncCollector<Signed<IAgentStep>> outputStream)
        {
            var state = await context.FetchStateAsync<State>();

            await Task.WhenAll(
                currentProposerStream.ForEachAsync(async currentProposer =>
                {
                    if (currentProposer.Equals(self))
                    {
                        outputStream.AddAsync(state.LastCommit);
                    }
                }, CancellationToken.None),

                committerStream.ForEachAsync(async commit =>
                {
                    state.LastCommit = commit;
                    await context.UpdateStateAsync(state);
                }, CancellationToken.None));
        }
    }
}