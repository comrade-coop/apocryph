using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks.Dataflow;
using Apocryph.Consensus;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Dataflow;
using Perper.WebJobs.Extensions.Model;
using Perper.WebJobs.Extensions.Triggers;

namespace Apocryph.Routing.FunctionApp
{
    using Apocryph.PerperUtilities;
    public static class RouterInput
    {
        [FunctionName("RouterInput")]
        public static IAsyncEnumerable<Message> RunAsync([PerperTrigger] (IAsyncEnumerable<Message> calls, IAsyncEnumerable<List<Reference>> subscriptions) input, IContext context, IStateEntry<List<Reference>?> lastSubscriptions)
        {
            var output = ApocryphDataflow.EmptyBlock<Message>();

            var subscriptions = ApocryphDataflow.KeepLastBlock(lastSubscriptions);
            input.subscriptions.ToDataflow().LinkTo(subscriptions);

            var subscriber = ApocryphDataflow.SubsciberBlock<Reference, Message>(async reference =>
            {
                var (targetInput, targetOutput) = await context.CallFunctionAsync<(string, IStream<Message>)>("GetChainInstance", reference.Chain);

                return targetOutput.Replay().ToDataflow(); // TODO: Make sure to replay only messages newer than the subscription
            });

            subscriptions.LinkTo(subscriber);
            subscriber.LinkTo(output);

            input.calls.ToDataflow().LinkTo(output);

            return output.ToAsyncEnumerable();
        }
    }
}