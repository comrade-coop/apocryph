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
    public static class IpfsOutput
    {
        [FunctionName(nameof(IpfsOutput))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("ipfsGateway")] string ipfsGateway,
            [Perper("topic")] string topic,
            [PerperStream("dataStream")] IAsyncEnumerable<ISigned<object>> dataStream,
            ILogger logger)
        {
            var ipfs = new IpfsClient(ipfsGateway);

            await dataStream.ForEachAsync(async item =>
            {
                try
                {
                    var bytes = Encoding.UTF8.GetBytes(JsonConvert.SerializeObject(item, typeof(ISigned<object>), IpfsJsonSettings.DefaultSettings));

                    await ipfs.PubSub.PublishAsync(topic, bytes, CancellationToken.None);
                }
                catch (Exception e)
                {
                    logger.LogError(e.ToString());
                }
            }, CancellationToken.None);
        }
    }
}