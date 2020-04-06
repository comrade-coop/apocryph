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
    public static class CallNotificationStepSplitter
    {
        [FunctionName(nameof(CallNotificationStepSplitter))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("notificationsStream")] IAsyncEnumerable<IHashed<CallNotification>> notificationsStream,
            [PerperStream("outputStream")] IAsyncCollector<Hash> outputStream)
        {
            await notificationsStream.ForEachAsync(async notification =>
            {
                await outputStream.AddAsync(notification.Value.Step);
            }, CancellationToken.None);
        }
    }
}