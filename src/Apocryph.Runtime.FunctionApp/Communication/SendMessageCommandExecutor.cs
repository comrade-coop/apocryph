using System.Collections.Generic;
using System.Linq;
using System.Security.Cryptography;
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
    public static class SendMessageCommandExecutor
    {
        public class State
        {
            public Dictionary<string, IHashed<ValidatorSet>> ValidatorSets { get; set; }
        }

        [FunctionName(nameof(SendMessageCommandExecutor))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("agentId")] string agentId,
            [Perper("ipfsGateway")] string ipfsGateway,
            [PerperStream("otherValidatorSetsStream")] IAsyncEnumerable<Dictionary<string, IHashed<ValidatorSet>>> otherValidatorSetsStream,
            [PerperStream("commandsStream")] IAsyncEnumerable<AgentCommand> commandsStream,
            CancellationToken cancellationToken)
        {
            var state = await context.FetchStateAsync<State>() ?? new State();

            await Task.WhenAll(
                otherValidatorSetsStream.ForEachAsync(async validatorSets =>
                {
                    state.ValidatorSets = validatorSets;
                    await context.UpdateStateAsync(state);
                }, cancellationToken),

                commandsStream.ForEachAsync(async command =>
                {
                    var notification = new CallNotification
                    {
                        From = agentId,
                        Command = command,
                        // Step = TODO,
                        // ValidatorSet = TODO,
                        // Commits = TODO
                    };

                    // await context.CallWorkerAsync<object>(nameof(PBFTNotificationWorker), new
                    // {
                    //     agentId = command.Receiver.Issuer,
                    //     validatorSet = state.ValidatorSets[command.Receiver.Issuer],
                    //     notification = notification,
                    //     ipfsGateway
                    // }, cancellationToken);
                }, cancellationToken));
        }
    }
}