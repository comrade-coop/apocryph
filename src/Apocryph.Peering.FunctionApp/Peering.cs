using System.Collections.Generic;
using System.Threading.Tasks;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Model;
using Perper.WebJobs.Extensions.Triggers;

namespace Apocryph.Peering.FunctionApp
{
    public class Peering
    {
        private IState _state;

        public Peering(IState state)
        {
            _state = state;
        }

        [FunctionName("Apocryph-Peering")]
        public void Start([PerperTrigger] object? input)
        {
        }

        [FunctionName("Connect")]
        public async Task<IAsyncEnumerable<object>> Connect([PerperTrigger] (Peer peer, string connectionType, IAsyncEnumerable<object> messages) input)
        {
            var key = input.peer.Id.ToString();
            var handler = await _state.GetValue<PeerHandler>(key, () => default!);

            return await handler.InvokeAsync((input.connectionType, input.messages));
        }

        [FunctionName("Register")]
        public Task Register([PerperTrigger] (Peer peer, PeerHandler handler) input)
        {
            var key = input.peer.Id.ToString();
            return _state.SetValue(key, input.handler);
        }
    }
}