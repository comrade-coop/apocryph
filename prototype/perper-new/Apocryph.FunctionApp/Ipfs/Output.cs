using System.Text;
using System.Threading;
using System.Threading.Tasks;
﻿using Ipfs.Http;
using Microsoft.Azure.WebJobs;
using Newtonsoft.Json;
using Perper.WebJobs.Extensions.Bindings;
using Perper.WebJobs.Extensions.Model;
using Perper.WebJobs.Extensions.Triggers;

namespace Apocryph.FunctionApp.Ipfs
{
    public static class Output
    {
        [FunctionName("IpfsOutput")]
        public static async Task Run([PerperStreamTrigger] IPerperStreamContext context,
            [PerperStream("ipfsAddress")] string ipfsAddress,
            [PerperStream("topic")] string topic,
            [PerperStream("data")] IPerperStream<object> dataStream)
        {
            var ipfs = new IpfsClient(ipfsAddress);

            await dataStream.Listen(async item => {
                var bytes = Encoding.UTF8.GetBytes(JsonConvert.SerializeObject(item));
                await ipfs.PubSub.PublishAsync(topic, bytes, CancellationToken.None);
            }, CancellationToken.None);
        }
    }
}