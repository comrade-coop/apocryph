using System.Text;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Model;
using Ipfs.Http;
using Microsoft.Azure.WebJobs;
using Newtonsoft.Json;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.FunctionApp.Ipfs
{
    public static class Input
    {
        [FunctionName("IpfsInput")]
        public static async Task Run([PerperTrigger("IpfsInput")] IPerperStreamContext context,
            [Perper("ipfsGateway")] string ipfsGateway,
            [Perper("topic")] string topic,
            [Perper("outputStream")] IAsyncCollector<Signed<object>> outputStream)
        {
            var ipfs = new IpfsClient(ipfsGateway);

            await ipfs.PubSub.SubscribeAsync(topic, async message => {
                var bytes = message.DataBytes;
                // FIXME: Do not blindly trust that Hash and Value match and that Signature, Hash, and Signer match
                var item = JsonConvert.DeserializeObject<Signed<object>>(Encoding.UTF8.GetString(bytes));
                await outputStream.AddAsync(item);
            }, CancellationToken.None);
        }
    }
}