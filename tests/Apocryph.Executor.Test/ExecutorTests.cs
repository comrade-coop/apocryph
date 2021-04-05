using System.Collections.Generic;
using System.Linq;
using Apocryph.Consensus;
using Apocryph.HashRegistry.Test;
using Xunit;

namespace Apocryph.Executor.Test
{
    public class ExecutorTests
    {
        [Fact]
        public async void TestAgentScenario_ProducesExpectedMessages()
        {
            var hashRegistry = HashRegistryFakes.GetHashRegistryProxy();
            var executor = await ExecutorFakes.GetExecutor(ExecutorFakes.TestAgents);
            var (chain, inputMessages, expectedOutputMessages) = await ExecutorFakes.GetTestAgentScenario(hashRegistry, "-", null, 1);

            var agentStates = await chain.GenesisState.AgentStates.EnumerateItems(hashRegistry).ToDictionaryAsync(x => x.Nonce, x => x);

            var outputMessages = new List<Message>();

            foreach (var inputMessage in inputMessages)
            {
                var inputState = agentStates[inputMessage.Target.AgentNonce];
                var (resultState, resultMessages) = await executor.InvokeAsync((inputState, inputMessage));

                Assert.Equal(inputState, resultState, SerializedComparer.Instance);

                outputMessages.AddRange(resultMessages);
            }

            Assert.Equal(outputMessages.ToArray(), expectedOutputMessages, SerializedComparer.Instance);

        }
    }
}