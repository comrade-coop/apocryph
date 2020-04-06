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
    public static class SubscriptionCommandExecutor
    {
        [FunctionName(nameof(SubscriptionCommandExecutor))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("ipfsGateway")] string ipfsGateway,
            [PerperStream("otherValidatorSetsStream")] IAsyncEnumerable<Dictionary<string, IHashed<ValidatorSet>>> otherValidatorSetsStream,
            [PerperStream("commandsStream")] IAsyncEnumerable<SubscriptionCommand> commandsStream,
            CancellationToken cancellationToken)
        {
            var cts = new CancellationTokenSource(); // Dispose!
            await using var outputStreams = new AsyncDisposableList();
            await using var utilityStreams = new AsyncDisposableList();
            await commandsStream.ForEachAsync(async subscription =>
            {
                var otherId = subscription.Target;
                var validatorSetsStream = await context.StreamFunctionAsync(nameof(AgentZeroValidatorSetsSplitter), new
                {
                    agentId = otherId,
                    otherValidatorSetsStream,
                });

                var lightNodeStream = await context.StreamFunctionAsync(nameof(PBFTLightNode), new
                {
                    agentId = otherId,
                    validatorSetsStream,
                    ipfsGateway
                });

                var commandsStream = await context.StreamFunctionAsync(nameof(CommandSplitter), new
                {
                    stepsStream = lightNodeStream
                });

                var outputStream = await context.StreamFunctionAsync(nameof(SubscriptionCommandOutput), new
                {
                    otherId,
                    publicationsStream = commandsStream
                });

                utilityStreams.Add(validatorSetsStream);
                utilityStreams.Add(lightNodeStream);
                utilityStreams.Add(commandsStream);
                outputStreams.Add(outputStream);

                cts.Dispose();
                cts = new CancellationTokenSource();
                Task.Run(() =>
                {
                    context.BindOutput(outputStreams.ToArray(), cts.Token);
                });
            }, cancellationToken);
        }
    }
}