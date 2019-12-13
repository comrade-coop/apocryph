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
    public static class Validator
    {
        private class State
        {
            public ValidatorKey Proposer { get; set; }
            public Hash CurrentStep { get; set; }
        }

        [FunctionName("Validator")]
        public static async Task Run([PerperStream("Validator")] IPerperStreamContext context,
            [Perper("validatorSet")] ValidatorSet validatorSet,
            [Perper("committerStream")] IAsyncEnumerable<Hash> committerStream,
            [Perper("currentProposerStream")] IAsyncEnumerable<ValidatorKey> currentProposerStream,
            [Perper("proposalsStream")] IAsyncEnumerable<Signed<IAgentStep>> proposalsStream,
            [Perper("outputStream")] IAsyncCollector<Hash> outputStream)
        {
            var state = context.GetState<State>();

            await Task.WhenAll(
                currentProposerStream.ForEachAsync(async currentProposer =>
                {
                    state.Proposer = currentProposer;

                    await context.SaveState();
                }, CancellationToken.None),

                committerStream.ForEachAsync(async commit =>
                {
                    state.CurrentStep = commit;

                    await context.SaveState();
                }, CancellationToken.None),

                proposalsStream.ForEachAsync(async proposal =>
                {
                    if (state.Proposer.Equals(proposal.Signer) && state.CurrentStep == proposal.Value.Previous)
                    {
                        await outputStream.AddAsync(proposal.Value.Previous);
                    }
                }, CancellationToken.None));
        }
    }
}