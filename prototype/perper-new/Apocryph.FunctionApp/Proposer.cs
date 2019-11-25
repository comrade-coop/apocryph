using System.Threading;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Model;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Bindings;
using Perper.WebJobs.Extensions.Model;
using Perper.WebJobs.Extensions.Triggers;

namespace Apocryph.FunctionApp
{
    public static class Proposer
    {
        [FunctionName("Proposer")]
        public static async Task Run([PerperStreamTrigger] IPerperStreamContext context,
            [PerperStream("proposalsStream")] IPerperStream<(AgentInput, AgentOutput)> proposalsStream,
            [PerperStream] IAsyncCollector<(AgentInput, AgentOutput)> outputStream)
        {
            await proposalsStream.Listen(async proposal => await outputStream.AddAsync(proposal), CancellationToken.None);
        }
    }
}