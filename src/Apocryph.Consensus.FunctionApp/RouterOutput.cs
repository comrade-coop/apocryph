using System.Collections.Generic;
using System.Threading.Tasks;
using Apocryph.Ipfs;
using Apocryph.ServiceRegistry;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Bindings;
using Perper.WebJobs.Extensions.Model;
using Perper.WebJobs.Extensions.Triggers;

namespace Apocryph.Consensus.FunctionApp
{
    public static class RouterOutput
    {
        [FunctionName("RouterOutput")]
        public static async IAsyncEnumerable<Message> RunAsync([PerperTrigger] (IAsyncEnumerable<Message> outbox, Hash<Chain> self, IAgent serviceRegistry) input, IContext context)
        {
            await foreach (var message in input.outbox)
            {
                if (message.Target.Chain == input.self) // FIXME: Needs a better way to handle publications, as this currently disallows self-messages
                {
                    yield return message;
                }
                else
                {
                    var locator = new ServiceLocator("Chain", message.Target.Chain.ToString());
                    var targetChain = await input.serviceRegistry.CallFunctionAsync<Service>("Lookup", locator);
                    var stream = targetChain.Inputs["messages"];

                    await context.CallActionAsync("PostMessage", (stream, message));
                }
            }
        }

        // HACK: IContext does not provide a way to post a message to a blank stream directly
        [FunctionName("PostMessage")]
        public static Task PostMessage([PerperTrigger(ParameterExpression = "{'stream': 0}")] (string _, Message message) input, [Perper(Stream = "{stream}")] IAsyncCollector<Message> collector)
        {
            return collector.AddAsync(input.message);
        }
    }
}