using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Agent;
using Apocryph.FunctionApp.Command;
using Apocryph.FunctionApp.Model;
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
        public class State
        {
            public Dictionary<string, ValidatorSet> ValidatorSets { get; set; }
        }

        [FunctionName("SubscriptionCommandExecutor")]
        public static async Task Run([PerperStreamTrigger("SubscriptionCommandExecutor")] IPerperStreamContext context,
            [Perper("ipfsGateway")] string ipfsGateway,
            [Perper("validatorSetsStream")] IAsyncEnumerable<Dictionary<string, ValidatorSet>> validatorSetsStream,
            [Perper("commandsStream")] IAsyncEnumerable<SubscriptionCommand> commandsStream,
            [Perper("outputStream")] IAsyncCollector<(string, object)> outputStream)
        {
            var ipfs = new IpfsClient(ipfsGateway);
            var state = await context.FetchStateAsync<State>();

            await Task.WhenAll(
                commandsStream.ForEachAsync(subscription =>
                {
                    var agentId = subscription.Target;
                    var topic = "apocryph-agent-" + agentId;
                    ipfs.PubSub.SubscribeAsync(topic, async message =>
                    {
                        var bytes = message.DataBytes;
                        // FIXME: Do not blindly trust that Hash and Value match and that Signature, Hash, and Signer match
                        var input = JsonConvert.DeserializeObject<Signed<AgentInput>>(Encoding.UTF8.GetString(bytes));
                        if (input != null && input.Value.CommitSignatures
                            .All(kv => kv.Key.ValidateSignature(input.Value.Previous, kv.Value)))
                        {
                            var validatorSet = state.ValidatorSets[agentId];
                            var committed = input.Value.CommitSignatures.Keys
                                .Select(signer => validatorSet.Weights[signer]).Sum();
                            if (3 * committed > 2 * validatorSet.Total)
                            {
                                var hash = input.Value.Previous;
                                // NOTE: Currently blocks other items on the stream and does not process them
                                // -- we should at least timeout
                                // FIXME: Should use DAG/IPLD API instead
                                var block = await ipfs.Block.GetAsync(Cid.Read(hash.Bytes), CancellationToken.None);

                                var output = JsonConvert.DeserializeObject<AgentOutput>(Encoding.UTF8.GetString(block.DataBytes));

                                foreach (var command in output.Commands)
                                {
                                    if (command is PublicationCommand publication)
                                    {
                                        await outputStream.AddAsync((agentId, publication.Payload));
                                    }
                                }
                            }
                        }
                    }, CancellationToken.None);
                }, CancellationToken.None),

                validatorSetsStream.ForEachAsync(async validatorSets =>
                {
                    state.ValidatorSets = validatorSets;
                    await context.UpdateStateAsync(state);
                }, CancellationToken.None));
        }
    }
}