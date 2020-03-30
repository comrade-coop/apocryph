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
        public class State
        {
            public Dictionary<string, IHashed<ValidatorSet>> ValidatorSets { get; set; }
        }
        [FunctionName(nameof(CallNotificationProcessor))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("agentId")] string agentId,
            [Perper("ipfsGateway")] string ipfsGateway,
            [PerperStream("validatorSetsStream")] IAsyncEnumerable<Dictionary<string, IHashed<ValidatorSet>>> validatorSetsStream,
            [PerperStream("notificationsStream")] IAsyncEnumerable<ISigned<CallNotification>> notificationsStream,
            CancellationToken cancellationToken)
        {
            var cts = new CancellationTokenSource(); // Dispose!
            await using var outputStreams = new AsyncDisposableList();
            await using var utilityStreams = new AsyncDisposableList();

            var state = await context.FetchStateAsync<State>() ?? new State();

            await Task.WhenAll(
                validatorSetsStream.ForEachAsync(async validatorSets =>
                {
                    state.ValidatorSets = validatorSets;
                    await context.UpdateStateAsync(state);
                }, cancellationToken),

                notificationsStream.ForEachAsync(async notification =>
                {
                    if (!state.ValidatorSets.ContainsKey(notification.Value.From) ||
                        !state.ValidatorSets[notification.Value.From].Value.Weights.ContainsKey(notification.Signer))
                    {
                        return;
                    }

                    var otherId = notification.Value.From;
                    var splitValidatorSetsStream = await context.StreamFunctionAsync(nameof(AgentZeroValidatorSetsSplitter), new
                    {
                        agentId = otherId,
                        validatorSetsStream,
                    });

                    var lightNodeStream = await context.StreamFunctionAsync(nameof(PBFTLightNode), new
                    {
                        agentId = otherId,
                        validatorSetsStream = splitValidatorSetsStream,
                        ipfsGateway
                    });

                    var commandsStream = await context.StreamFunctionAsync(nameof(CommandSplitter), new
                    {
                        stepsStream = lightNodeStream
                    });

                    var outputStream = await context.StreamFunctionAsync(nameof(CallNotificationOutput), new
                    {
                        agentId,
                        otherId,
                        commandsStream = commandsStream
                    });

                    utilityStreams.Add(splitValidatorSetsStream);
                    utilityStreams.Add(lightNodeStream);
                    utilityStreams.Add(commandsStream);
                    outputStreams.Add(outputStream);

                    cts.Dispose();
                    cts = new CancellationTokenSource();
                    Task.Run(() =>
                    {
                        context.BindOutput(outputStreams.ToArray(), cts.Token);
                    });
                }, cancellationToken));
        }
    }
}