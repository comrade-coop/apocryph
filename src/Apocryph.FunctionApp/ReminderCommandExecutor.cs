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
    public static class ReminderCommandExecutor
    {
        [FunctionName("ReminderCommandExecutor")]
        public static async Task Run([PerperStreamTrigger("ReminderCommandExecutor")] IPerperStreamContext context,
            [Perper("commandsStream")] IAsyncEnumerable<ReminderCommand> commandsStream,
            [Perper("outputStream")] IAsyncCollector<(string, object)> outputStream)
        {
            commandsStream.ForEachAsync(async reminder =>
            {
                Task.Run(async () => {
                    await Task.Delay(reminder.Time);
                    outputStream.AddAsync(("", reminder.Data));
                });
            }, CancellationToken.None);
        }
    }
}