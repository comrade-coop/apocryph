using System.Collections.Generic;
using System.Linq;
using System.Threading.Channels;
using Apocryph.Consensus.FunctionApp;
using Apocryph.HashRegistry;
using Apocryph.HashRegistry.Test;
using Apocryph.ServiceRegistry;
using Apocryph.ServiceRegistry.Test;
using Perper.WebJobs.Extensions.Fake;
using Perper.WebJobs.Extensions.Model;
using Xunit;
using Xunit.Abstractions;

namespace Apocryph.Consensus.Test
{
    public class RouterTests
    {
        private readonly ITestOutputHelper _output;
        public RouterTests(ITestOutputHelper output)
        {
            _output = output;
        }

        private Message GenerateMessage(int chain, int i)
        {
            var hash = Hash.From(chain);
            return new Message(new Reference(hash.Cast<Chain>(), 0, new[] { typeof(int).FullName! }), ReferenceData.From(i));
        }

        [Fact]
        public async void RouterInput_SubscribesAndPasses_AllMessages()
        {
            var hashRegistry = HashRegistryFakes.GetHashRegistryProxy();
            var (serviceRegistryAgent, serviceRegistry) = ServiceRegistryFakes.GetServiceRegistryAgent();

            var callMessages = new[] { GenerateMessage(0, 0), GenerateMessage(0, 1) };

            var subscriptionMessages = new[] { GenerateMessage(1, 2), GenerateMessage(1, 3) };
            var subscriptionsChannel = Channel.CreateUnbounded<List<Reference>>();

            var routerStateEntry = await ((IState)new FakeState()).Entry<List<Reference>?>("-", () => null);

            var routedMessages = RouterInput.RunAsync((callMessages.ToAsyncEnumerable(), subscriptionsChannel.Reader.ReadAllAsync(), serviceRegistryAgent), routerStateEntry);

            var routedMessagesEnumerator = routedMessages.GetAsyncEnumerator();

            foreach (var expectedMessage in callMessages)
            {
                Assert.True(await routedMessagesEnumerator.MoveNextAsync());
                Assert.Equal(routedMessagesEnumerator.Current, expectedMessage, SerializedComparer.Instance);
            }

            var subscriptionChain = Hash.From("-").Cast<Chain>();
            var subscriptionReference = new Reference(subscriptionChain, 0, new string[] { });

            await serviceRegistry.Register((new ServiceLocator("Chain", subscriptionChain.ToString()), new Service(new Dictionary<string, string>(), new Dictionary<string, IStream>()
            {
                {"messages", new FakeStream<Message>(subscriptionMessages)}
            })), default);

            await subscriptionsChannel.Writer.WriteAsync(new List<Reference> { subscriptionReference });

            foreach (var expectedMessage in subscriptionMessages)
            {
                Assert.True(await routedMessagesEnumerator.MoveNextAsync());
                Assert.Equal(routedMessagesEnumerator.Current, expectedMessage, SerializedComparer.Instance);
            }

            subscriptionsChannel.Writer.Complete();
        }

        [Fact]
        public async void RouterOutput_PostsAndPublishes_AllMessages()
        {
            var hashRegistry = HashRegistryFakes.GetHashRegistryProxy();
            var (serviceRegistryAgent, serviceRegistry) = ServiceRegistryFakes.GetServiceRegistryAgent();

            var messagesToPublish = new[] { GenerateMessage(0, 0), GenerateMessage(0, 1) };
            var messagesToSend = new[] { GenerateMessage(1, 2), GenerateMessage(1, 3) };
            var selfChainId = messagesToPublish[0].Target.Chain;

            var outputtedMessages = new List<(string, Message)>();

            foreach (var chainId in messagesToSend.Select(x => x.Target.Chain).Where(id => id != selfChainId).Distinct())
            {
                await serviceRegistry.Register((new ServiceLocator("Chain", chainId.ToString()), new Service(new Dictionary<string, string>(){
                    {"messages", chainId.ToString() + "-stream"}
                }, new Dictionary<string, IStream>())), default);
            }

            var context = new FakeContext();
            context.Agent.RegisterFunction("PostMessage", ((string, Message) input) => outputtedMessages.Add(input));

            var publishedMessages = RouterOutput.RunAsync((messagesToPublish.Concat(messagesToSend).ToAsyncEnumerable(), selfChainId, serviceRegistryAgent), context);

            Assert.Equal(await publishedMessages.ToArrayAsync(), messagesToPublish, SerializedComparer.Instance);
            Assert.Equal(outputtedMessages.ToArray(), messagesToSend.Select(x => (x.Target.Chain.ToString() + "-stream", x)), SerializedComparer.Instance);
        }
    }
}