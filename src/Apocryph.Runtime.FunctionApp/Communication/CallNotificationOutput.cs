using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Agent;
using Apocryph.Runtime.FunctionApp.Consensus;
using Apocryph.Runtime.FunctionApp.Ipfs;
using Ipfs;
using Ipfs.Http;
using Microsoft.Azure.WebJobs;
using Newtonsoft.Json;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;
using Block = Apocryph.Runtime.FunctionApp.Consensus.Core.Block;

namespace Apocryph.Runtime.FunctionApp.Communication
{
    public static class CallNotificationOutput
    {
        [FunctionName(nameof(CallNotificationOutput))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("agentId")] string agentId,
            [PerperStream("notificationsStream")] IAsyncEnumerable<IHashed<CallNotification>> notificationsStream,
            [PerperStream("notificationStepSplitterStream")] IAsyncEnumerable<IHashed<Block>> notificationStepSplitterStream,
            [PerperStream("outputStream")] IAsyncCollector<(string, object)> outputStream)
        {
            await using var stepSplitterStreamEnumerator = notificationStepSplitterStream.GetAsyncEnumerator();
            await notificationsStream.ForEachAsync(async notification =>
            {
                await stepSplitterStreamEnumerator.MoveNextAsync();
                var step = stepSplitterStreamEnumerator.Current;

                // TODO: Use Merkle proofs for this
                var found = false;

                foreach (var command in step.Value.Commands)
                {
                    if (command == notification.Value.Command)
                    {
                        found = true;
                    }
                }

                if (found)
                {
                    await outputStream.AddAsync((notification.Value.From, notification.Value.Command.Message));
                }
            }, CancellationToken.None);
        }
    }
}