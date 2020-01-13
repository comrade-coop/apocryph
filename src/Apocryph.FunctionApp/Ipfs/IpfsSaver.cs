using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Model;
﻿using Ipfs.Http;
using Microsoft.Azure.WebJobs;
using Microsoft.Extensions.Logging;
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
            [PerperStream("outputStream")] IAsyncCollector<IHashed<object>> outputStream,
            ILogger logger)
        {
            var ipfs = new IpfsClient(ipfsGateway);

            await dataStream.ForEachAsync(async item => {
                try
                {
                    var json = JsonConvert.SerializeObject(item, typeof(ISigned<object>), IpfsJsonSettings.DefaultSettings);
                    var bytes = Encoding.UTF8.GetBytes(json);

                    // FIXME: Should use DAG/IPLD API instead
                    var cid = await ipfs.Block.PutAsync(bytes, cancel: CancellationToken.None);

                    var hash = new Hash {Bytes = cid.ToArray()};

                    logger.LogDebug("Saved {json} as {hash} in ipfs", json, hash);

                    var hashedType = typeof(Hashed<>).MakeGenericType(item.GetType());
                    var hashed = (IHashed<object>)Activator.CreateInstance(hashedType, item, hash);

                    await outputStream.AddAsync(hashed);
                }
                catch (Exception e)
                {
                    logger.LogError(e.ToString());
                }
            }, CancellationToken.None);
        }
    }
}