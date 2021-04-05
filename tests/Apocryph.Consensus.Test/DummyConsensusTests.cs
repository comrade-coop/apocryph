using System.Linq;
using Apocryph.Consensus.Dummy.FunctionApp;
using Apocryph.Executor.Test;
using Apocryph.HashRegistry.Test;
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
            var executor = await ExecutorFakes.GetExecutor(ExecutorFakes.TestAgents);
            var (chain, inputMessages, expectedOutputMessages) = await ExecutorFakes.GetTestAgentScenario(hashRegistry, "Apocryph-DummyConsensus", null, 1);

            var dummyConsensus = new DummyConsensus();

            var outputMessages = await dummyConsensus.ExecutionStream((inputMessages.ToAsyncEnumerable(), hashRegistry, chain, executor)).ToArrayAsync();

            Assert.Equal(outputMessages, expectedOutputMessages, SerializedComparer.Instance);
        }
    }
}