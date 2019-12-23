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
    public static class ValidatorFilter
    {
        private class State
        {
            public ValidatorKey Proposer { get; set; }
            public Hash CurrentStep { get; set; }
        }

        [FunctionName("ValidatorFilter")]
        public static async Task Run([PerperStreamTrigger("ValidatorFilter")] IPerperStreamContext context,
            [Perper("committerStream")] IAsyncEnumerable<Hash> committerStream,
            [Perper("currentProposerStream")] IAsyncEnumerable<ValidatorKey> currentProposerStream,
            [Perper("proposalsStream")] IAsyncEnumerable<Signed<IAgentStep>> proposalsStream,
            [Perper("outputStream")] IAsyncCollector<Signed<IAgentStep>> outputStream)
        {
            var state = await context.FetchStateAsync<State>();

            await Task.WhenAll(
                currentProposerStream.ForEachAsync(async currentProposer =>
                {
                    state.Proposer = currentProposer;

                    await context.UpdateStateAsync(state);
                }, CancellationToken.None),

                committerStream.ForEachAsync(async commit =>
                {
                    state.CurrentStep = commit;

                    await context.UpdateStateAsync(state);
                }, CancellationToken.None),

                proposalsStream.ForEachAsync(async proposal =>
                {
                    if (state.Proposer.Equals(proposal.Signer) && state.CurrentStep == proposal.Value.Previous)
                    {
                        await outputStream.AddAsync(proposal);
                    }
                }, CancellationToken.None));
        }
    }
}