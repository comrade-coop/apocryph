using System.Collections.Concurrent;
using System.Text.Json;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Ipfs.Serialization;

namespace Apocryph.Ipfs.Fake
{
    public class FakeHashResolver : IHashResolver
    {
        public ConcurrentDictionary<Hash<object>, TaskCompletionSource<byte[]>> Values { get; set; } = new ConcurrentDictionary<Hash<object>, TaskCompletionSource<byte[]>>();

        public Task<Hash<T>> StoreAsync<T>(T value, CancellationToken cancellationToken = default)
        {
            var serialized = JsonSerializer.SerializeToUtf8Bytes(value, ApocryphSerializationOptions.JsonSerializerOptions);
            var hash = Hash.FromBytes<object>(serialized);

            var taskSource = Values.GetOrAdd(hash, _ => new TaskCompletionSource<byte[]>());
            taskSource.TrySetResult(serialized);

            return Task.FromResult(hash.Cast<T>());
        }

        public async Task<T> RetrieveAsync<T>(Hash<T> hash, CancellationToken cancellationToken = default)
        {
            var taskSource = Values.GetOrAdd(hash.Cast<object>(), _ => new TaskCompletionSource<byte[]>());

            var serialized = await taskSource.Task.ContinueWith(x => x.Result, cancellationToken);
            var value = JsonSerializer.Deserialize<T>(serialized, ApocryphSerializationOptions.JsonSerializerOptions);

            return value;
        }
    }
}