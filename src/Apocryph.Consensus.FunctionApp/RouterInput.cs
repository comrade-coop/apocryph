using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks.Dataflow;
using Apocryph.ServiceRegistry;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Dataflow;
using Perper.WebJobs.Extensions.Model;
using Perper.WebJobs.Extensions.Triggers;

namespace Apocryph.Consensus.FunctionApp
{
    using Apocryph.PerperUtilities;
    public static class RouterInput
    {
        [FunctionName("RouterInput")]
        public static IAsyncEnumerable<Message> RunAsync([PerperTrigger] (IAsyncEnumerable<Message> calls, IAsyncEnumerable<List<Reference>> subscriptions, IAgent serviceRegistry) input, IStateEntry<List<Reference>?> lastSubscriptions)
        {
            var output = ApocryphDataflow.EmptyBlock<Message>();

            var subscriptions = ApocryphDataflow.KeepLastBlock(lastSubscriptions);
            input.subscriptions.ToDataflow().LinkTo(subscriptions);

            var subscriber = ApocryphDataflow.SubsciberBlock<Reference, Message>(async reference =>
            {
                var locator = new ServiceLocator("Chain", reference.Chain.ToString());
                var targetChain = await input.serviceRegistry.CallFunctionAsync<Service>("Lookup", locator);
                var stream = (IStream<Message>)targetChain.Outputs["messages"];

                return stream.Replay().ToDataflow(); // TODO: Make sure to replay only messages newer than the subscription
            });

            subscriptions.LinkTo(subscriber);
            subscriber.LinkTo(output);

            input.calls.ToDataflow().LinkTo(output);

            return output.ToAsyncEnumerable();
        }
    }
}