using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using System.Threading.Tasks.Dataflow;
using Apocryph.Consensus;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Dataflow;
using Perper.WebJobs.Extensions.Model;
using Perper.WebJobs.Extensions.Triggers;

namespace Apocryph.Routing.FunctionApp
{
    public static class RouterInput
    {
        [FunctionName("RouterInput")]
        public static IAsyncEnumerable<Message> RunAsync([PerperTrigger] (IAsyncEnumerable<Message> calls, IAsyncEnumerable<List<Reference>> subscriptions) input, IContext context, IStateEntry<List<Reference>?> lastSubscriptions)
        {
            var output = EmptyBlock<Message>();

            var subscriptions = KeepLastBlock(lastSubscriptions);
            input.subscriptions.ToDataflow().LinkTo(subscriptions);

            var subscriber = SubsciberBlock<Reference, Message>(async reference =>
            {
                var (targetInput, targetOutput) = await context.CallFunctionAsync<(string, IStream<Message>)>("GetChainInstance", reference.Chain);

                return targetOutput.Replay().ToDataflow(); // TODO: Make sure to replay only messages newer than the subscription
            });

            subscriptions.LinkTo(subscriber);
            subscriber.LinkTo(output);

            input.calls.ToDataflow().LinkTo(output);

            return output.ToAsyncEnumerable();
        }

        private static IPropagatorBlock<T, T> EmptyBlock<T>()
        {
            return new BufferBlock<T>(new DataflowBlockOptions { BoundedCapacity = 1 });
        }

        private static IPropagatorBlock<T, T> KeepLastBlock<T>(IStateEntry<T?> stateEntry) where T : class, new()
        {
            var output = new TransformBlock<T, T>(value =>
            {
                stateEntry.Value = value;
                return value;
            });

            output.Post(stateEntry.Value ?? new T());

            return output;
        }

        private static IPropagatorBlock<IEnumerable<TKey>, TValue> SubsciberBlock<TKey, TValue>(Func<TKey, Task<ISourceBlock<TValue>>> resolver)
            where TKey : notnull
        {
            var links = new Dictionary<TKey, IDisposable>();
            var output = EmptyBlock<TValue>();
            var subscriber = new ActionBlock<IEnumerable<TKey>>(async (subscriptions) =>
            {
                var seenSubscriptions = new HashSet<TKey>();
                foreach (var subscription in subscriptions)
                {
                    if (!links.ContainsKey(subscription))
                    {
                        var source = await resolver(subscription);
                        links[subscription] = source.LinkTo(output);
                    }
                    seenSubscriptions.Add(subscription);
                }

                foreach (var (subscription, link) in links)
                {
                    if (!seenSubscriptions.Contains(subscription))
                    {
                        link.Dispose();
                    }
                }
            });

            return DataflowBlock.Encapsulate(subscriber, output);
        }
    }
}