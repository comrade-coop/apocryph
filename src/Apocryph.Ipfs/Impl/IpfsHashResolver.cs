using System.Text.Json;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Ipfs.Serialization;
using Ipfs;
using Ipfs.Http;

namespace Apocryph.Ipfs.Impl
{
    public class IpfsHashResolver : IHashResolver
    {
        private readonly IpfsClient _ipfs;

        public IpfsHashResolver(IpfsClient ipfs)
        {
            _ipfs = ipfs;
        }

        public async Task<Hash<T>> StoreAsync<T>(T value, CancellationToken cancellationToken = default)
        {
            var serialized = JsonSerializer.SerializeToUtf8Bytes(value, ApocryphSerializationOptions.JsonSerializerOptions);

            // NOTE: Using "raw" here instead of "json", since Ipfs.Http.Client doesn't seem to consider "json" a valid MultiCodec
            await _ipfs.Block.PutAsync(serialized, "raw", "sha2-256", "identity", false);

            return Hash.FromBytes<T>(serialized);
        }

        public async Task<T> RetrieveAsync<T>(Hash<T> hash, CancellationToken cancellationToken = default)
        {
            var cid = new Cid { ContentType = "raw", Hash = new MultiHash("sha2-256", hash.Bytes) };

            // NOTE: The Ipfs.Http.Client library uses a GET request for Block.GetAsync, which doesn't work since go-ipfs v0.5.
            // See https://github.com/richardschneider/net-ipfs-http-client/issues/62 for more details.
            // var block = await _ipfs.Block.GetAsync(multihash);

            var stream = await _ipfs.PostDownloadAsync("block/get", cancellationToken, cid);

            return (await JsonSerializer.DeserializeAsync<T>(stream, ApocryphSerializationOptions.JsonSerializerOptions, cancellationToken))!;
        }
    }
}