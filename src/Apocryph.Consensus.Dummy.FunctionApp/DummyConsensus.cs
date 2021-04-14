using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Apocryph.Ipfs;
using Apocryph.KoTH;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Model;
using Perper.WebJobs.Extensions.Triggers;

namespace Apocryph.Consensus.Dummy.FunctionApp
{
    public class DummyConsensus
    {
        [FunctionName("Apocryph-DummyConsensus")]
        public async Task<IAsyncEnumerable<Message>> Start([PerperTrigger] (
            IAsyncEnumerable<Message> messages,
            string subscriptionsStream,
            Chain chain,
            IAsyncEnumerable<(Hash<Chain>, Slot?[])> kothStates,
            IAgent executor) input, IContext context)
        {
            return await context.StreamFunctionAsync<Message>("ExecutionStream", (input.messages, input.chain, input.executor));
        }

        [FunctionName("ExecutionStream")]
        public async IAsyncEnumerable<Message> ExecutionStream([PerperTrigger] (IAsyncEnumerable<Message> messages, Chain chain, IAgent executor) input, IHashResolver hashResolver)
        {
            var agentStates = await input.chain.GenesisState.AgentStates.EnumerateItems(hashResolver).ToDictionaryAsync(x => x.Nonce, x => x);
            var self = Hash.From(input.chain);

            await foreach (var message in input.messages)
            {
                if (message.Target.Chain != self || !message.Target.AllowedMessageTypes.Contains(message.Data.Type))
                    continue;

                var state = agentStates[message.Target.AgentNonce];

                var (newState, resultMessages) = await input.executor.CallFunctionAsync<(AgentState, Message[])>("Execute", (self, state, message));

                agentStates[message.Target.AgentNonce] = newState;

                foreach (var resultMessage in resultMessages)
                {
                    yield return resultMessage;
                }
            }
        }
    }
}