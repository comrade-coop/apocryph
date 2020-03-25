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
        [FunctionName(nameof(ReminderCommandExecutor))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("commandsStream")] IAsyncEnumerable<ReminderCommand> commandsStream,
            [PerperStream("outputStream")] IAsyncCollector<(string, object)> outputStream,
            CancellationToken cancellationToken)
        {
            await commandsStream.ForEachAsync(reminder =>
            {
                Task.Run(async () => {
                    await Task.Delay(reminder.Time, cancellationToken);
                    await outputStream.AddAsync(("Reminder", reminder.Data));
                });
            }, cancellationToken);
        }
    }
}