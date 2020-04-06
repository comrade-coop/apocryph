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
    public static class CallNotificationProcessor
    {
        [FunctionName(nameof(CallNotificationProcessor))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("agentId")] string agentId,
            [Perper("ipfsGateway")] string ipfsGateway,
            [Perper("otherValidatorSetsStream")] object[] otherValidatorSetsStream,
            [Perper("notificationsStream")] object[] notificationsStream,
            CancellationToken cancellationToken)
        {
            await using var notificationValidatorStream = await context.StreamFunctionAsync(nameof(CallNotificationValidator), new
            {
                agentId,
                otherValidatorSetsStream,
                notificationsStream
            });

            await using var _notificationStepSplitterStream = await context.StreamFunctionAsync(nameof(CallNotificationStepSplitter), new
            {
                notificationsStream = notificationValidatorStream,
            });

            await using var notificationStepSplitterStream = await context.StreamFunctionAsync(nameof(IpfsLoader), new
            {
                ipfsGateway,
                hashStream = _notificationStepSplitterStream,
            });

            await using var outputStream = await context.StreamFunctionAsync(nameof(CallNotificationOutput), new
            {
                agentId,
                notificationsStream = notificationValidatorStream,
                notificationStepSplitterStream
            });

            await context.BindOutput(outputStream, cancellationToken);
        }
    }
}