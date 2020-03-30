using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Agent;
using Apocryph.FunctionApp.Command;
using Apocryph.FunctionApp.Model;
using Apocryph.FunctionApp.Ipfs;
using Ipfs;
using Ipfs.Http;
using Microsoft.Azure.WebJobs;
using Newtonsoft.Json;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.FunctionApp
{
    public static class CallNotificationOutput
    {
        [FunctionName(nameof(CallNotificationOutput))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("otherId")] string otherId,
            [Perper("agentId")] string agentId,
            [PerperStream("commandsStream")] IAsyncEnumerable<SendMessageCommand> commandsStream,
            [PerperStream("outputStream")] IAsyncCollector<(string, object)> outputStream)
        {
            await commandsStream.ForEachAsync(async command =>
            {
                if (command.Target == agentId)
                {
                    await outputStream.AddAsync((otherId, command.Payload));
                }
            }, CancellationToken.None);
        }
    }
}