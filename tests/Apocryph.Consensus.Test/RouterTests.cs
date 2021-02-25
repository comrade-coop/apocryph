using System.Collections.Generic;
using System.Linq;
using System.Threading.Channels;
using Apocryph.Consensus.FunctionApp;
using Apocryph.HashRegistry;
using Apocryph.ServiceRegistry;
using Perper.WebJobs.Extensions.Fake;
using Perper.WebJobs.Extensions.Model;
using Xunit;
using Xunit.Abstractions;

namespace Apocryph.Consensus.Test
{
    using HashRegistry = Apocryph.HashRegistry.FunctionApp.HashRegistry;
    using ServiceRegistry = Apocryph.ServiceRegistry.FunctionApp.ServiceRegistry;

    public class RouterTests
    {
        private readonly ITestOutputHelper _output;
        public RouterTests(ITestOutputHelper output)
        {
            _output = output;
        }

        private HashRegistryProxy GetHashRegistryProxy()
        {
            var registry = new HashRegistry(new FakeState());

            var agent = new FakeAgent();
            agent.RegisterFunction("Store", (byte[] data) => registry.Store(data, default));
            agent.RegisterFunction("Retrieve", (Hash hash) => registry.Retrieve(hash, default));

            return new HashRegistryProxy(agent);
        }

        private (FakeAgent, ServiceRegistry) GetServiceRegistryAgent()
        {
            var registry = new ServiceRegistry(new FakeState());

            var agent = new FakeAgent();
            agent.RegisterFunction("Register", ((ServiceLocator locator, Service service) input) => registry.Register(input, default));
            agent.RegisterFunction("RegisterTypeHandler", ((string type, Handler handler) input) => registry.RegisterTypeHandler(input, default));
            agent.RegisterFunction("Lookup", (ServiceLocator input) => registry.Lookup(input, default));

            return (agent, registry);
        }

        private Message GenerateMessage(int chain, int i)
        {
            var hash = Hash.From(chain);
            return new Message(new Reference(hash.Cast<Chain>(), hash.Cast<AgentState>(), new[] { typeof(int).FullName! }), ReferenceData.From(i));
        }

        [Fact]
        public async void RouterInput_SubscribesAndPasses_AllMessages()
        {
            var hashRegistry = GetHashRegistryProxy();
            var (serviceRegistryAgent, serviceRegistry) = GetServiceRegistryAgent();

            var callMessages = new[] { GenerateMessage(0, 0), GenerateMessage(0, 1) };

            var subscriptionMessages = new[] { GenerateMessage(1, 2), GenerateMessage(1, 3) };
            var subscriptionsChannel = Channel.CreateUnbounded<List<Reference>>();

            var routerStateEntry = await ((IState)new FakeState()).Entry<List<Reference>?>("-", () => null);

            var routerInput = new RouterInput();
            var routedMessages = routerInput.RunAsync((callMessages.ToAsyncEnumerable(), subscriptionsChannel.Reader.ReadAllAsync(), serviceRegistryAgent), routerStateEntry);

            var routedMessagesEnumerator = routedMessages.GetAsyncEnumerator();

            foreach (var expectedMessage in callMessages)
            {
                Assert.True(await routedMessagesEnumerator.MoveNextAsync());
                Assert.Equal(routedMessagesEnumerator.Current, expectedMessage, SerializedComparer.Instance);
            }

            var subscriptionChain = Hash.From("-").Cast<Chain>();
            var subscriptionReference = new Reference(subscriptionChain, subscriptionChain.Cast<AgentState>(), new string[] { });

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
            var hashRegistry = GetHashRegistryProxy();
            var (serviceRegistryAgent, serviceRegistry) = GetServiceRegistryAgent();

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

            var routerOutput = new RouterOutput();
            var publishedMessages = routerOutput.RunAsync((messagesToPublish.Concat(messagesToSend).ToAsyncEnumerable(), selfChainId, serviceRegistryAgent), context);

            Assert.Equal(await publishedMessages.ToArrayAsync(), messagesToPublish, SerializedComparer.Instance);
            Assert.Equal(outputtedMessages.ToArray(), messagesToSend.Select(x => (x.Target.Chain.ToString() + "-stream", x)), SerializedComparer.Instance);
        }
    }
}