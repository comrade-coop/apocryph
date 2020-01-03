using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Model;
﻿using Ipfs.Http;
using Microsoft.Azure.WebJobs;
using Newtonsoft.Json;
using Perper.WebJobs.Extensions.Bindings;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;
using Perper.WebJobs.Extensions.Triggers;

namespace Apocryph.FunctionApp.Ipfs
{
    public static class IpfsSaver
    {
        [FunctionName(nameof(IpfsSaver))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("ipfsGateway")] string ipfsGateway,
            [PerperStream("dataStream")] IAsyncEnumerable<object> dataStream,
            [PerperStream("outputStream")] IAsyncCollector<Hashed<object>> outputStream)
        {
            var ipfs = new IpfsClient(ipfsGateway);

            await dataStream.ForEachAsync(async item => {
                var bytes = Encoding.UTF8.GetBytes(JsonConvert.SerializeObject(item));

                // FIXME: Should use DAG/IPLD API instead
                var cid = await ipfs.Block.PutAsync(bytes, cancel: CancellationToken.None);

                var hash = new Hash {Bytes = cid.ToArray()};

                await outputStream.AddAsync(new Hashed<object>(item, hash));
            }, CancellationToken.None);
        }
    }
}