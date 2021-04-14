using System.Collections.Generic;
using System.Linq;
using Apocryph.Consensus;
using Apocryph.Ipfs.Fake;
using Apocryph.Ipfs.Test;
using Xunit;

namespace Apocryph.Executor.Test
{
    public class ExecutorTests
    {
        [Fact]
        public async void TestAgentScenario_ProducesExpectedMessages()
        {
            var hashResolver = new FakeHashResolver();
            var executor = await ExecutorFakes.GetExecutor(ExecutorFakes.TestAgents);
            var (chain, inputMessages, expectedOutputMessages) = await ExecutorFakes.GetTestAgentScenario(hashResolver, "-", null, 1);
            var chainId = await hashResolver.StoreAsync(chain);

            var agentStates = await chain.GenesisState.AgentStates.EnumerateItems(hashResolver).ToDictionaryAsync(x => x.Nonce, x => x);

            var outputMessages = new List<Message>();

            foreach (var inputMessage in inputMessages)
            {
                var inputState = agentStates[inputMessage.Target.AgentNonce];
                var (resultState, resultMessages) = await executor.CallFunctionAsync<(AgentState, Message[])>("Execute", (chainId, inputState, inputMessage));

                Assert.Equal(inputState, resultState, SerializedComparer.Instance);

                outputMessages.AddRange(resultMessages);
            }

            Assert.Equal(outputMessages.ToArray(), expectedOutputMessages, SerializedComparer.Instance);

        }
    }
}