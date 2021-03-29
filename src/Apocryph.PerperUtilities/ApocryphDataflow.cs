using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using System.Threading.Tasks.Dataflow;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.PerperUtilities
{
    public static class ApocryphDataflow
    {
        public static IPropagatorBlock<T, T> EmptyBlock<T>()
        {
            return new BufferBlock<T>(new DataflowBlockOptions { BoundedCapacity = 1 });
        }

        public static IPropagatorBlock<T, T> KeepLastBlock<T>(IStateEntry<T?> stateEntry) where T : class, new()
        {
            var output = new TransformBlock<T, T>(value =>
            {
                stateEntry.Value = value;
                return value;
            });

            output.Post(stateEntry.Value ?? new T());

            return output;
        }

        public static IPropagatorBlock<IEnumerable<TKey>, TValue> SubsciberBlock<TKey, TValue>(Func<TKey, Task<ISourceBlock<TValue>>> resolver)
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

        // FIXME: obsoleted by https://github.com/obecto/perper/commit/c21139721cce283a9544619830b67ca8bb5fbee6
        public static ISourceBlock<T> ToDataflow<T>(this IAsyncEnumerable<T> enumerable, CancellationToken cancellationToken = default)
        {
            var block = new BufferBlock<T>(new DataflowBlockOptions { CancellationToken = cancellationToken, BoundedCapacity = 1 });

            async Task helper()
            {
                await foreach (var item in enumerable.WithCancellation(cancellationToken))
                {
                    await block.SendAsync(item);
                }
            }

            helper().ContinueWith(completedTask =>
            {
                if (completedTask.Status == TaskStatus.Faulted) ((IDataflowBlock)block).Fault(completedTask.Exception!);
                else if (completedTask.Status == TaskStatus.RanToCompletion) block.Complete();
            });

            return block;
        }
    }
}