using System.Linq;
using Apocryph.Consensus.Dummy.FunctionApp;
using Apocryph.HashRegistry;
using Apocryph.HashRegistry.MerkleTree;
using Apocryph.HashRegistry.Test;
using Perper.WebJobs.Extensions.Fake;
using Xunit;
using Xunit.Abstractions;

namespace Apocryph.Consensus.Test
{
    public class DummyConsensusTests
    {
        private readonly ITestOutputHelper _output;
        public DummyConsensusTests(ITestOutputHelper output)
        {
            _output = output;
        }

        [Fact]
        public async void ExecutionStream_Returns_ExpectedMesages()
        {
            var hashRegistry = HashRegistryFakes.GetHashRegistryProxy();

            var agentStates = new[] {
                new AgentState(ReferenceData.From("123"), "Agent1"),
                new AgentState(ReferenceData.From(null), "Agent2")
            };
            var agentStatesTree = (await MerkleTreeBuilder.CreateFromValues(hashRegistry, agentStates, 2)).First().GetRoot();
            var chain = new Chain(agentStatesTree, "Apocryph-DummyConsensus", 1);

            var chainId = Hash.From(chain);
            var agentIds = agentStates.Select(x => Hash.From(x)).ToArray();
            var testMessageAllowed = new string[] { typeof(int).FullName! };

            var inputMessages = new Message[]
            {
                new Message(new Reference(chainId, agentIds[0], testMessageAllowed), ReferenceData.From(4)),
                new Message(new Reference(chainId, agentIds[1], testMessageAllowed), ReferenceData.From(3)),
            };

            var expectedOutputMessages = new Message[]
            {
                new Message(new Reference(chainId, agentIds[1], testMessageAllowed), ReferenceData.From(3)),
                new Message(new Reference(chainId, agentIds[0], testMessageAllowed), ReferenceData.From(2)),
            };

            var context = new FakeContext();
            context.Agent.RegisterAgent("Agent1", () => context.Agent);
            context.Agent.RegisterFunction("Agent1", ((AgentState state, Message message) input) =>
            {
                Assert.Equal(input.state, agentStates[0], SerializedComparer.Instance);
                var result = input.message.Data.Deserialize<int>() - 1;
                return (input.state, new[] { new Message(new Reference(chainId, agentIds[1], testMessageAllowed), ReferenceData.From(result)) });
            });
            context.Agent.RegisterAgent("Agent2", () => context.Agent);
            context.Agent.RegisterFunction("Agent2", ((AgentState state, Message message) input) =>
            {
                Assert.Equal(input.state, agentStates[1], SerializedComparer.Instance);
                var result = input.message.Data.Deserialize<int>() - 1;
                return (input.state, new[] { new Message(new Reference(chainId, agentIds[0], testMessageAllowed), ReferenceData.From(result)) });
            });

            var dummyConsensus = new DummyConsensus(context);

            var outputMessages = await dummyConsensus.ExecutionStream((inputMessages.ToAsyncEnumerable(), hashRegistry, chain)).ToArrayAsync();

            Assert.Equal(outputMessages, expectedOutputMessages, SerializedComparer.Instance);
        }
    }
}