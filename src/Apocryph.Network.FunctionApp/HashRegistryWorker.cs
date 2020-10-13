using System.Threading;
using System.Threading.Tasks;
using System.Text.Json;
using Apocryph.Core.Consensus.Blocks;
using Block = Apocryph.Core.Consensus.Blocks.Block;
using Apocryph.Core.Consensus.Serialization;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;
using Ipfs.Http;
using Ipfs;

namespace Apocryph.Runtime.FunctionApp
{
    public class HashRegistryWorker
    {
        private static IpfsClient? _ipfsClient;

        [FunctionName(nameof(HashRegistryWorker))]
        [return: Perper("$return")]
        public async Task<Block> Run([PerperWorkerTrigger] PerperWorkerContext context,
            [Perper("hash")] Hash hash)
        {
            if (_ipfsClient == null)
            {
                _ipfsClient = new IpfsClient();
            }

            var cid = new Cid { ContentType = "raw", Hash = new MultiHash("sha2-256", hash.Value) };

            // FIXME: The Ipfs.Http.Client library uses a GET request for Block.GetAsync, which doesn't work since go-ipfs v0.5.
            // See https://github.com/richardschneider/net-ipfs-http-client/issues/62 for more details.
            // var block = await _ipfsClient.Block.GetAsync(multihash);

            var stream = await _ipfsClient.PostDownloadAsync("block/get", default(CancellationToken), cid);

            return await JsonSerializer.DeserializeAsync<Block>(stream, ApocryphSerializationOptions.JsonSerializerOptions);
        }
    }
}