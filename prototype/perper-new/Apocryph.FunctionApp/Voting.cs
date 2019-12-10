using System.Collections.Generic;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Model;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Bindings;
using Perper.WebJobs.Extensions.Model;
using Perper.WebJobs.Extensions.Triggers;

namespace Apocryph.FunctionApp
{
    public static class Voting
    {
        [FunctionName("Voting")]
        public static async Task Run([PerperStreamTrigger] IPerperStreamContext context,
            [PerperStream("runtimeStream")] IAsyncEnumerable<(IAgentStep, bool)> runtimeStream,
            [PerperStream("proposalsStream")] IAsyncEnumerable<IAgentStep> proposalsStream,
            [PerperStream] IAsyncCollector<object> outputStream)
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