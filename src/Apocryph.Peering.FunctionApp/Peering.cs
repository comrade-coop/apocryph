using System.Collections.Generic;
using System.Threading.Tasks;
using Apocryph.PerperUtilities;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Model;
using Perper.WebJobs.Extensions.Triggers;

namespace Apocryph.Peering.FunctionApp
{
    using PeerHandler = Handler<IAsyncEnumerable<object>>;

    public class Peering
    {
        public class PeeringState
        {
            public Peer Self;
            public IAgent Router;

            public PeeringState(Peer self, IAgent router)
            {
                Self = self;
                Router = router;
            }
        }

        private IState _state;

        public Peering(IState state)
        {
            _state = state;
        }

        [FunctionName("_GetHandler")]
        public async Task<PeerHandler> _GetHandler([PerperTrigger] (Peer peer, string connectionType) input)
        {
            var key = $"{input.peer}-{input.connectionType}";
            var handler = await _state.GetValue<PeerHandler>(key, () => default!);

            return handler;
        }

        [FunctionName("_GetPeering")]
        public async Task<IAgent> _GetPeering([PerperTrigger] Peer peer, IContext context)
        {
            var (agent, _) = await context.StartAgentAsync<object?>("Apocryph-Peering", new PeeringState(peer, context.Agent));

            return agent;
        }

        [FunctionName("_SetHandler")]
        public Task _SetHandler([PerperTrigger] (Peer peer, string connectionType, PeerHandler handler) input)
        {
            var key = $"{input.peer}-{input.connectionType}";
            return _state.SetValue<PeerHandler>(key, input.handler);
        }

        [FunctionName("Apocryph-Peering")]
        public async Task Start([PerperTrigger] PeeringState? input)
        {
            if (input != null)
            {
                await _state.SetValue("state", input);
            }
        }


        [FunctionName("Connect")]
        public async Task<IAsyncEnumerable<object>> Connect([PerperTrigger] (Peer peer, string connectionType, IAsyncEnumerable<object> messages) input)
        {
            var state = await _state.GetValue<PeeringState>("state", () => default!);
            var handler = await state.Router.CallFunctionAsync<PeerHandler>("_GetHandler", (input.peer, input.connectionType));

            return await handler.InvokeAsync((state.Self, input.messages));
        }

        [FunctionName("Register")]
        public async Task Register([PerperTrigger] (string connectionType, PeerHandler handler) input)
        {
            var state = await _state.GetValue<PeeringState>("state", () => default!);

            await state.Router.CallActionAsync("_SetHandler", (state.Self, input.connectionType, input.handler));
        }
    }
}