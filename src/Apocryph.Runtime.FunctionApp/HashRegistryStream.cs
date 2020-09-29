using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using System.Text.Json;
using Apocryph.Core.Consensus.Blocks;
using Apocryph.Core.Consensus.Serialization;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;
using Ipfs.Http;
using Ipfs;

namespace Apocryph.Runtime.FunctionApp
{
    public class HashRegistryStream
    {
        [FunctionName(nameof(HashRegistryStream))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("filter")] Type filter,
            [Perper("input")] IAsyncEnumerable<object> input,
            [Perper("output")] IAsyncCollector<HashRegistryEntry> output)
        {
            var ipfs = new IpfsClient();

            try
            {
                await foreach (var value in input)
                {
                    if (filter.IsAssignableFrom(value.GetType()))
                    {
                        var serialized = JsonSerializer.SerializeToUtf8Bytes(value, ApocryphSerializationOptions.JsonSerializerOptions);

                        // FIXME: Using "raw" here instead of "json", since Ipfs.Http.Client doesn't seem to consider "json" a valid MultiCodec
                        var cid = await ipfs.Block.PutAsync(serialized, "raw", "sha2-256");
                    }
                }
            }
            catch (Exception e)
            {
                Console.WriteLine(e);
            }
        }

        private static IpfsClient? _ipfsClient;

        public static async Task<T?> GetObjectByHash<T>(IQueryable<HashRegistryEntry> queryable, Hash hash)
            where T: class
        {
            if (_ipfsClient == null)
            {
                _ipfsClient = new IpfsClient();
            }

            var cid = new Cid {ContentType = "raw", Hash = new MultiHash("sha2-256", hash.Value)};

            // FIXME: The Ipfs.Http.Client library uses a GET request for Block.GetAsync, which doesn't work since go-ipfs v0.5.
            // See https://github.com/richardschneider/net-ipfs-http-client/issues/62 for more details.

            // var block = await _ipfsClient.Block.GetAsync(multihash);

            var stream = await _ipfsClient.PostDownloadAsync("block/get", default(CancellationToken), cid);
            return await JsonSerializer.DeserializeAsync<T?>(stream, ApocryphSerializationOptions.JsonSerializerOptions);
        }
    }
}