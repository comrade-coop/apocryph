using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Runtime.FunctionApp.Ipfs;
﻿using Ipfs.Http;
using Microsoft.Azure.WebJobs;
using Microsoft.Extensions.Logging;
using Newtonsoft.Json;
using Perper.WebJobs.Extensions.Bindings;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;
using Perper.WebJobs.Extensions.Triggers;

namespace Apocryph.Runtime.FunctionApp.Ipfs
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

            await dataStream.ForEachAsync(async item =>
            {
                try
                {
                    var jToken = IpfsJsonSettings.JTokenFromObject(item);
                    var cid = await ipfs.Dag.PutAsync(jToken, cancel: CancellationToken.None);
                    var hash = new Cid {Bytes = cid.ToArray()};

                    logger.LogDebug("Saved {json} as {hash} in ipfs", jToken.ToString(Formatting.None), hash);

                    await outputStream.AddAsync(Hashed.Create(item, hash));
                }
                catch (Exception e)
                {
                    logger.LogError(e.ToString());
                }
            }, CancellationToken.None);
        }
    }
}