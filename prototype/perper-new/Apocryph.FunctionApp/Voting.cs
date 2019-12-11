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

        [FunctionName("Voting")]
        public static async Task Run([PerperTrigger("Voting")] IPerperStreamContext context,
            [Perper("runtimeStream")] IAsyncEnumerable<(Hashed<IAgentStep>, bool)> runtimeStream,
            [Perper("proposalsStream")] IAsyncEnumerable<Signed<IAgentStep>> proposalsStream,
            [Perper("outputStream")] IAsyncCollector<Vote> outputStream)
        {
            var state = context.GetState<State>();

            await Task.WhenAll(
                proposalsStream.ForEachAsync(async proposal =>
                {
                    state.ExpectedNextSteps[proposal.Value.Previous] = proposal.Hash;

                    await context.SaveState();
                }, CancellationToken.None),

                runtimeStream.ForEachAsync(async item =>
                {
                    var (nextStep, isProposal) = item;
                    if (!isProposal && state.ExpectedNextSteps[nextStep.Value.Previous] == nextStep.Hash)
                    {
                        await outputStream.AddAsync(new Vote { For = nextStep.Hash });
                    }
                }, CancellationToken.None));
        }
    }
}