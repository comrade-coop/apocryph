using System.Collections.Generic;
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

        [FunctionName("Voting")]
        public static async Task Run([Perper(Stream = "Voting")] IPerperStreamContext context,
            [Perper("runtimeStream")] IAsyncEnumerable<(IAgentStep, bool)> runtimeStream,
            [Perper("proposalsStream")] IAsyncEnumerable<IAgentStep> proposalsStream,
            [Perper("outputStream")] IAsyncCollector<object> outputStream)
        {
            var state = context.GetState<State>("state");

            await Task.WhenAll(
                proposalsStream.Listen(async proposal =>
                {
                    state.ExpectedNextSteps[proposal.PreviousHash] = proposal.Hash;

                    await context.SaveState("state", state);
                }, CancellationToken.None),

                runtimeStream.Listen(async item =>
                {
                    var (nextStep, isProposal) = item;
                    if (!isProposal && state.ExpectedNextSteps[nextStep.PreviousHash] == nextStep.Hash)
                    {
                        await outputStream.AddAsync(new Vote { ForHash = nextStep.Hash });
                    }
                }, CancellationToken.None));
        }
    }
}