using System.Threading;
using System.Threading.Tasks;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Model;
using Perper.WebJobs.Extensions.Triggers;

namespace Apocryph.HashRegistry.FunctionApp
{
    public class HashRegistry
    {
        private IState _state;

        public HashRegistry(IState state)
        {
            _state = state;
        }

        [FunctionName("Store")]
        public async Task Store([PerperTrigger] byte[] serialized, CancellationToken cancellationToken)
        {
            var hash = Hash.FromSerialized<object?>(serialized);
            await _state.SetValue(hash.ToString(), serialized);
        }

        [FunctionName("Retrieve")]
        public async Task<byte[]> Retrieve([PerperTrigger] Hash hash, CancellationToken cancellationToken)
        {
            byte[]? serialized = null;
            while (serialized == null)
            {
                cancellationToken.ThrowIfCancellationRequested();

                serialized = await _state.GetValue<byte[]?>(hash.ToString(), () => null);
                await Task.Delay(50, cancellationToken); // Simulate IPFS retying to get the hash
            }
            return serialized;
        }
    }
}