using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Model;
using Microsoft.Azure.WebJobs;
using Microsoft.Extensions.Logging;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.FunctionApp
{
    public static class ValidatorFilter
    {
        private class State
        {
            public State(bool init = false)
            {
                if (init)
                {
                    CurrentStep = new Hash {Bytes = new byte[] {0}};
                }
            }

            public ValidatorKey Proposer { get; set; }
            public Hash CurrentStep { get; set; }
        }

        [FunctionName(nameof(ValidatorFilter))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("committerStream")] IAsyncEnumerable<IHashed<IAgentStep>> committerStream,
            [PerperStream("currentProposerStream")] IAsyncEnumerable<ValidatorKey> currentProposerStream,
            [PerperStream("proposalsStream")] IAsyncEnumerable<ISigned<IAgentStep>> proposalsStream,
            [PerperStream("outputStream")] IAsyncCollector<ISigned<IAgentStep>> outputStream,
            ILogger logger)
        {
            var state = await context.FetchStateAsync<State>() ?? new State(true);

            await Task.WhenAll(
                currentProposerStream.ForEachAsync(async currentProposer =>
                {
                    state.Proposer = currentProposer;

                    await context.UpdateStateAsync(state);
                }, CancellationToken.None),

                committerStream.ForEachAsync(async commit =>
                {
                    state.CurrentStep = commit.Hash;

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