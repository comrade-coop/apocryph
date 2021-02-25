using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Channels;
using System.Threading.Tasks;
using Apocryph.ServiceRegistry;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Model;
using Perper.WebJobs.Extensions.Triggers;

namespace Apocryph.Consensus.FunctionApp
{
    public class RouterInput
    {
        private List<Task> _tasks = new List<Task>();
        private Dictionary<Reference, CancellationTokenSource> _cancellationTokens = new Dictionary<Reference, CancellationTokenSource>();
        private Channel<Message> _channel = Channel.CreateUnbounded<Message>();
        private IAgent? _serviceRegistry = null;

        [FunctionName("RouterInput")]
        public IAsyncEnumerable<Message> RunAsync([PerperTrigger] (IAsyncEnumerable<Message> calls, IAsyncEnumerable<List<Reference>> subscriptions, IAgent serviceRegistry) input, IStateEntry<List<Reference>?> lastSubscriptions)
        {
            _serviceRegistry = input.serviceRegistry;

            _tasks.Add(UpdateSubscriptions(input.subscriptions, lastSubscriptions));
            _tasks.Add(IterateStream(input.calls, default));

            return _channel.Reader.ReadAllAsync();
        }

        private async Task UpdateSubscriptions(IAsyncEnumerable<List<Reference>> subscriptions, IStateEntry<List<Reference>?> lastSubscriptions)
        {
            if (lastSubscriptions.Value == null)
            {
                lastSubscriptions.Value = new List<Reference>();
            }

            UpdateSubscriptions(lastSubscriptions.Value!);
            await foreach (var newSubscriptions in subscriptions)
            {
                lastSubscriptions.Value = newSubscriptions;
                UpdateSubscriptions(lastSubscriptions.Value!);
            }
        }

        private void UpdateSubscriptions(List<Reference> subscriptions)
        {
            var seenSubscriptions = new HashSet<Reference>();
            foreach (var subscription in subscriptions)
            {
                if (!_cancellationTokens.ContainsKey(subscription))
                {
                    var tokenSource = new CancellationTokenSource();
                    _tasks.Add(IterateStream(subscription, tokenSource.Token));
                    _cancellationTokens[subscription] = tokenSource;
                }
                seenSubscriptions.Add(subscription);
            }

            foreach (var (oldSubscription, tokenSource) in _cancellationTokens)
            {
                if (!seenSubscriptions.Contains(oldSubscription))
                {
                    tokenSource.Cancel();
                }
            }
        }

        private async Task IterateStream(Reference subscription, CancellationToken token)
        {
            var locator = new ServiceLocator("Chain", subscription.Chain.ToString());
            var targetChain = await _serviceRegistry!.CallFunctionAsync<Service>("Lookup", locator);
            var stream = (IStream<Message>)targetChain.Outputs["messages"];

            await IterateStream(stream.Replay(), token); // TODO: Make sure to replay only newer messages somehow
        }

        private async Task IterateStream(IAsyncEnumerable<Message> stream, CancellationToken token)
        {
            await foreach (var message in stream.WithCancellation(token))
            {
                await _channel.Writer.WriteAsync(message);
            }
        }
    }
}