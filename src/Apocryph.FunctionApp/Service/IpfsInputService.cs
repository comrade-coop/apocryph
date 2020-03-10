using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Model;
using Apocryph.FunctionApp.Ipfs;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.FunctionApp.Service
{
    public static class IpfsInputService
    {
        [FunctionName("ServiceIpfsInput")]
        public static async Task RunLauncher([PerperStreamTrigger] PerperStreamContext context,
            [Perper("agentId")] string agentId,
            [Perper("ipfsGateway")] string ipfsGateway,
            [Perper("commandsStream")] object[] commandsStream,
            CancellationToken cancellationToken)
        {
            var topic = "apocryph-agentInput-" + agentId;

            await using var ipfsStream = await context.StreamFunctionAsync(nameof(IpfsInput), new
            {
                ipfsGateway,
                topic
            });

            await using var transformedStream = await context.StreamFunctionAsync(nameof(IpfsInputSubscriber), new
            {
                ipfsStream,
                commandsStream
            });

            await context.BindOutput(transformedStream, cancellationToken);
        }
    }
}