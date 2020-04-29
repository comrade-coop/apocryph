using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Agent;
using Apocryph.Runtime.FunctionApp.Consensus;
using Apocryph.Runtime.FunctionApp.Ipfs;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Runtime.FunctionApp.Communication
{
    public static class CommandSplitter
    {
        [FunctionName(nameof(CommandSplitter))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("stepsStream")] IAsyncEnumerable<IHashed<AgentBlock>> stepsStream,
            [PerperStream("outputStream")] IAsyncCollector<AgentCommand> outputStream)
        {
            await stepsStream.ForEachAsync(async output =>
            {
                foreach (var command in output.Value.Commands)
                {
                    await outputStream.AddAsync(command);
                }
            }, CancellationToken.None);
        }
    }
}