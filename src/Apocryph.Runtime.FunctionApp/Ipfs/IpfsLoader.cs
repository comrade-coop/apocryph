using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Runtime.FunctionApp.Ipfs;
﻿using Ipfs;
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
    public static class IpfsLoader
    {
        [FunctionName(nameof(IpfsLoader))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("ipfsGateway")] string ipfsGateway,
            [PerperStream("hashStream")] IAsyncEnumerable<Cid> hashStream,
            [PerperStream("outputStream")] IAsyncCollector<IHashed<object>> outputStream,
            ILogger logger)
        {
            var ipfs = new IpfsClient(ipfsGateway);

            await hashStream.ForEachAsync(async hash =>
            {
                try
                {
                    // NOTE: Currently blocks other items on the stream and does not process them
                    // -- we should at least timeout

                    var jToken = await ipfs.Dag.GetAsync(Cid.Read(hash.Bytes), CancellationToken.None);
                    var item = jToken.ToObject<object>(JsonSerializer.Create(IpfsJsonSettings.DefaultSettings));

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