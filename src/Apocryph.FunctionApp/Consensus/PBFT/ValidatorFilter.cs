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
            public ValidatorKey Proposer { get; set; }
        }

        [FunctionName(nameof(ValidatorFilter))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("currentProposerStream")] IAsyncEnumerable<ValidatorKey> currentProposerStream,
            [PerperStream("proposalsStream")] IAsyncEnumerable<ISigned<Proposal>> proposalsStream,
            [PerperStream("outputStream")] IAsyncCollector<Hash> outputStream,
            ILogger logger)
        {
            var state = await context.FetchStateAsync<State>() ?? new State();

            await Task.WhenAll(
                currentProposerStream.ForEachAsync(async currentProposer =>
                {
                    state.Proposer = currentProposer;

                    await context.UpdateStateAsync(state);
                }, CancellationToken.None),

                proposalsStream.ForEachAsync(async proposal =>
                {
                    if (state.Proposer.Equals(proposal.Signer))
                    {
                        await outputStream.AddAsync(proposal.Value.For);
                    }
                }, CancellationToken.None));
        }
    }
}