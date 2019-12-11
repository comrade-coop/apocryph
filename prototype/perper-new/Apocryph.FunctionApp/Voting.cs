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
        [FunctionName("Voting")]
        public static async Task Run([Perper(Stream = "Voting")] IPerperStreamContext context,
            [Perper("runtimeStream")] IAsyncEnumerable<(IAgentStep, bool)> runtimeStream,
            [Perper("proposalsStream")] IAsyncEnumerable<IAgentStep> proposalsStream,
            [Perper("outputStream")] IAsyncCollector<object> outputStream)
        {
            var expectedNextSteps = new Dictionary<IAgentStep, IAgentStep>();
            await Task.WhenAll(
                proposalsStream.Listen(proposal =>
                {
                    expectedNextSteps[proposal.Previous] = proposal;
                }, CancellationToken.None),

                runtimeStream.Listen(async item =>
                {
                    var (nextStep, isProposal) = item;
                    if (!isProposal && expectedNextSteps[nextStep.Previous] == nextStep)
                    {
                        await outputStream.AddAsync(new Vote { For = nextStep });
                    }
                }, CancellationToken.None));
        }
    }
}