using System.Diagnostics;
using System.Text.Json;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.HashRegistry.Serialization;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.HashRegistry
{
    public struct HashRegistryProxy
    {
        public IAgent HashRegistryAgent;

        public HashRegistryProxy(IAgent hashRegistryAgent)
        {
            HashRegistryAgent = hashRegistryAgent;
        }

        public async Task<Hash<T>> StoreAsync<T>(T value, CancellationToken cancellationToken = default)
        {
            var serialized = JsonSerializer.SerializeToUtf8Bytes(value, ApocryphSerializationOptions.JsonSerializerOptions);
            await HashRegistryAgent.CallActionAsync("Store", serialized); // cancellationToken
            var hash = Hash.FromSerialized<T>(serialized);
            return hash;
        }

        public async Task<T> RetrieveAsync<T>(Hash<T> hash, CancellationToken cancellationToken = default)
        {
            var serialized = await HashRegistryAgent.CallFunctionAsync<byte[]>("Retrieve", hash.Cast<object>()); // cancellationToken
            var value = JsonSerializer.Deserialize<T>(serialized, ApocryphSerializationOptions.JsonSerializerOptions);
            return value;
        }
    }
}