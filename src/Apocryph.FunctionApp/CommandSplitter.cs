using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Agent;
using Apocryph.FunctionApp.Command;
using Apocryph.FunctionApp.Model;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.FunctionApp
{
    public static class CommandSplitter
    {
        [FunctionName("CommandSplitter")]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("committerStream")] IAsyncEnumerable<Hashed<AgentOutput>> committerStream,
            [PerperStream("outputStream")] IAsyncCollector<ICommand> outputStream)
        {
            await committerStream.ForEachAsync(async output =>
            {
                foreach (var command in output.Value.Commands)
                {
                    outputStream.AddAsync(command);
                }
            }, CancellationToken.None);
        }
    }
}