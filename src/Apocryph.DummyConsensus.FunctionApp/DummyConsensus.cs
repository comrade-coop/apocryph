using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Apocryph.Consensus;
using Apocryph.HashRegistry;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Model;
using Perper.WebJobs.Extensions.Triggers;

namespace Apocryph.DummyConsensus.FunctionApp
{
    public class DummyConsensus
    {
        private IContext _context;

        public DummyConsensus(IContext context)
        {
            _context = context;
        }

        [FunctionName("Apocryph-DummyConsensus")]
        public async Task<IAsyncEnumerable<Message>> Start([PerperTrigger] (IAsyncEnumerable<Message> messages, string subscribtionsStream, HashRegistryProxy hashRegsitry, Chain chain) input)
        {
            return await _context.StreamFunctionAsync<Message>("ExecutionStream", (input.messages, input.hashRegsitry, input.chain));
        }

        [FunctionName("ExecutionStream")]
        public async IAsyncEnumerable<Message> ExecutionStream([PerperTrigger] (IAsyncEnumerable<Message> messages, HashRegistryProxy hashRegsitry, Chain chain) input)
        {
            var genesisAgentStates = await input.chain.GenesisStates.EnumerateItems(input.hashRegsitry).ToListAsync();
            var statesById = genesisAgentStates.ToDictionary(x => Hash.From(x), x => x);
            var self = Hash.From(input.chain);

            await foreach (var message in input.messages)
            {
                if (message.Target.Chain != self)
                    continue;

                if (!message.Target.AllowedMessageTypes.Contains(message.Data.Type))
                    continue;

                if (!statesById.ContainsKey(message.Target.Agent))
                    continue;

                Message[] resultMessages;
                (statesById[message.Target.Agent], resultMessages) = await Execute(statesById[message.Target.Agent], message);

                foreach (var resultMessage in resultMessages)
                {
                    yield return resultMessage;
                }
            }
        }

        private async Task<(AgentState, Message[])> Execute(AgentState agentState, Message message)
        {
            var (_, result) = await _context.StartAgentAsync<(AgentState, Message[])>(agentState.Handler, (agentState, message));

            return result;
        }
    }
}