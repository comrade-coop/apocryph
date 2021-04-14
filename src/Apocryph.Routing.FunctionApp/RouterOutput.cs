using System.Collections.Generic;
using System.Threading.Tasks;
using Apocryph.Consensus;
using Apocryph.Ipfs;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Bindings;
using Perper.WebJobs.Extensions.Model;
using Perper.WebJobs.Extensions.Triggers;

namespace Apocryph.Routing.FunctionApp
{
    public static class RouterOutput
    {
        [FunctionName("RouterOutput")]
        public static async Task RunAsync(
            [PerperTrigger(ParameterExpression = "{'stream': 0}")] (string _, IAsyncEnumerable<Message> outbox, Hash<Chain> self) input,
            [Perper(Stream = "{stream}")] IAsyncCollector<Message> collector,
            IContext context,
            IState state)
        {
            await foreach (var message in input.outbox)
            {
                if (message.Target.Chain == input.self && message.Target.AgentNonce < 0) // FIXME: Needs a better way to handle publications
                {
                    await collector.AddAsync(message);
                }
                else
                {
                    var (targetInput, targetOutput) = await context.CallFunctionAsync<(string, IStream<Message>)>("GetChainInstance", message.Target.Chain);

                    await context.CallActionAsync("PostMessage", (targetInput, message));
                }
            }
        }

        // HACK: IContext does not provide a way to post a message to a blank stream directly
        [FunctionName("PostMessage")]
        public static Task PostMessage(
                [PerperTrigger(ParameterExpression = "{'stream': 0}")] (string _, Message message) input,
                [Perper(Stream = "{stream}")] IAsyncCollector<Message> collector)
        {
            return collector.AddAsync(input.message);
        }
    }
}