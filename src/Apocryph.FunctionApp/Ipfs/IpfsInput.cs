using System;
using System.Text;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Model;
using Ipfs.Http;
using Microsoft.Azure.WebJobs;
using Microsoft.Extensions.Logging;
using Newtonsoft.Json;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.FunctionApp.Ipfs
{
    public static class IpfsInput
    {
        [FunctionName(nameof(IpfsInput))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("ipfsGateway")] string ipfsGateway,
            [Perper("topic")] string topic,
            [PerperStream("outputStream")] IAsyncCollector<ISigned<object>> outputStream,
            ILogger logger)
        {
            var ipfs = new IpfsClient(ipfsGateway);

            await ipfs.PubSub.SubscribeAsync(topic, async message =>
            {
                try
                {
                    var bytes = message.DataBytes;
                    // FIXME: Do not blindly trust that Hash and Value match and that Signature, Hash, and Signer match
                    var item = JsonConvert.DeserializeObject<ISigned<object>>(Encoding.UTF8.GetString(bytes), IpfsJsonSettings.DefaultSettings);

                    logger.LogTrace("Received a '{type}' from '{topic}' on IPFS pubsub", item.GetType(), topic);

                    await outputStream.AddAsync(item);
                }
                catch (Exception e)
                {
                    logger.LogError(e.ToString());
                }
            }, CancellationToken.None);

            await context.BindOutput(CancellationToken.None);
        }
    }
}