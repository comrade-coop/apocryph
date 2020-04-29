using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Agent;
using Apocryph.Runtime.FunctionApp.Ipfs;
using Ipfs;
using Ipfs.Http;
using Microsoft.Azure.WebJobs;
using Newtonsoft.Json;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Runtime.FunctionApp.Communication
{
    public static class CallNotificationStepSplitter
    {
        [FunctionName(nameof(CallNotificationStepSplitter))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("notificationsStream")] IAsyncEnumerable<IHashed<CallNotification>> notificationsStream,
            [PerperStream("outputStream")] IAsyncCollector<Cid> outputStream)
        {
            await notificationsStream.ForEachAsync(async notification =>
            {
                await outputStream.AddAsync(notification.Value.Block);
            }, CancellationToken.None);
        }
    }
}