using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Apocryph.HashRegistry;
using Apocryph.PerperUtilities;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Model;
using Perper.WebJobs.Extensions.Triggers;

namespace Apocryph.Consensus.Dummy.FunctionApp
{
    public class DummyConsensus
    {
        [FunctionName("Apocryph-DummyConsensus")]
        public async Task<IAsyncEnumerable<Message>> Start([PerperTrigger] (IAsyncEnumerable<Message> messages, string subscribtionsStream, HashRegistryProxy hashRegsitry, Chain chain, IHandler<(AgentState, Message[])> executor) input, IContext context)
        {
            return await context.StreamFunctionAsync<Message>("ExecutionStream", (input.messages, input.hashRegsitry, input.chain, input.executor));
        }

        [FunctionName("ExecutionStream")]
        public async IAsyncEnumerable<Message> ExecutionStream([PerperTrigger] (IAsyncEnumerable<Message> messages, HashRegistryProxy hashRegsitry, Chain chain, IHandler<(AgentState, Message[])> executor) input)
        {
            var genesisAgentStates = await input.chain.GenesisState.AgentStates.EnumerateItems(input.hashRegsitry).ToListAsync();
            var statesById = genesisAgentStates.ToDictionary(x => x.Nonce, x => x);
            var self = Hash.From(input.chain);

            await foreach (var message in input.messages)
            {
                if (message.Target.Chain != self)
                    continue;

                if (!message.Target.AllowedMessageTypes.Contains(message.Data.Type))
                    continue;

                if (!statesById.ContainsKey(message.Target.AgentNonce))
                    continue;

                Message[] resultMessages;
                (statesById[message.Target.AgentNonce], resultMessages) = await input.executor.InvokeAsync((statesById[message.Target.AgentNonce], message));

                foreach (var resultMessage in resultMessages)
                {
                    yield return resultMessage;
                }
            }
        }
    }
}