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
    public static class Input
    {
        [FunctionName("IpfsInput")]
        public static async Task Run([PerperStreamTrigger] IPerperStreamContext context,
            [PerperStream("ipfsAddress")] string ipfsAddress,
            [PerperStream("topic")] string topic,
            [PerperStream] IAsyncCollector<object> outputStream)
        {
            var ipfs = new IpfsClient(ipfsAddress);

            await ipfs.PubSub.SubscribeAsync(topic, async message => {
                var bytes = message.DataBytes;
                var item = JsonConvert.DeserializeObject(Encoding.UTF8.GetString(bytes));
                await outputStream.AddAsync(item);
            }, CancellationToken.None);
        }
    }
}