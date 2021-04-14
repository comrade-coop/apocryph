using System.Threading.Tasks;
using Apocryph.Consensus;
using Apocryph.Ipfs;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Model;
using Perper.WebJobs.Extensions.Triggers;

namespace Apocryph.Executor.FunctionApp
{
    public class Executor
    {
        private IState _state;

        public Executor(IState state)
        {
            _state = state;
        }

        [FunctionName("Apocryph-Executor")]
        public void Start([PerperTrigger] object? input)
        {
        }

        [FunctionName("Register")]
        public Task Register([PerperTrigger] (Hash<string> agentCodeHash, IAgent handlerAgent, string handlerFunction) input)
        {
            var key = $"{input.agentCodeHash}";
            return _state.SetValue(key, (input.handlerAgent, input.handlerFunction));
        }

        [FunctionName("Execute")]
        public async Task<(AgentState, Message[])> Execute([PerperTrigger] (Hash<Chain> chain, AgentState agent, Message message) input)
        {
            var key = $"{input.agent.CodeHash}";
            var (handlerAgent, handlerFunction) = await _state.GetValue<(IAgent, string)>(key, () => default!);
            return await handlerAgent.CallFunctionAsync<(AgentState, Message[])>(handlerFunction, (input.chain, input.agent, input.message));
        }
    }
}