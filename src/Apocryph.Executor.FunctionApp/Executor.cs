using System.Threading.Tasks;
using Apocryph.Consensus;
using Apocryph.HashRegistry;
using Apocryph.PerperUtilities;
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

        [FunctionName("_Register")]
        public Task _Register([PerperTrigger] (Hash<string> agentCodeHash, IHandler<(AgentState, Message[])> handler) input)
        {
            var key = $"{input.agentCodeHash}";
            return _state.SetValue(key, input.handler);
        }

        [FunctionName("Execute")]
        public async Task<(AgentState, Message[])> Execute([PerperTrigger] (AgentState agent, Message message) input)
        {
            var key = $"{input.agent.CodeHash}";
            var handler = await _state.GetValue<IHandler<(AgentState, Message[])>>(key, () => default!);
            return await handler.InvokeAsync((input.agent, input.message));
        }
    }
}