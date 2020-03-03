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
    public static class ServiceCommandFilter
    {
        [FunctionName(nameof(ServiceCommandFilter))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("commandsStream")] IAsyncEnumerable<ServiceCommand> commandsStream,
            [Perper("serviceName")] string serviceName,
            [PerperStream("outputStream")] IAsyncCollector<object> outputStream,
            CancellationToken cancellationToken)
        {
            await commandsStream.ForEachAsync(async command =>
            {
                if (command.Service == serviceName) {
                    await outputStream.AddAsync(command.Parameters);
                }
            }, cancellationToken);
        }
    }
}