using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Agent;
using Apocryph.Runtime.FunctionApp.Ipfs;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Runtime.FunctionApp.Communication
{
    public static class ReminderCommandExecutor
    {
        [FunctionName(nameof(ReminderCommandExecutor))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("commandsStream")] IAsyncEnumerable<AgentCommand> commandsStream,
            [PerperStream("outputStream")] IAsyncCollector<(string, object)> outputStream,
            CancellationToken cancellationToken)
        {
            await commandsStream.ForEachAsync(reminder =>
            {
                Task.Run(async () => {
                    await Task.Delay(reminder.Timeout, cancellationToken);
                    await outputStream.AddAsync(("Reminder", reminder.Message));
                });
            }, cancellationToken);
        }
    }
}