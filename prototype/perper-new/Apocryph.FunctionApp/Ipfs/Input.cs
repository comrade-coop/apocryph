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
        public static async Task Run([Perper(Stream="IpfsInput")] IPerperStreamContext context,
            [Perper("ipfsGateway")] string ipfsGateway,
            [Perper("topic")] string topic,
            [Perper("output")] IAsyncCollector<ISigned> output)
        {
            var ipfs = new IpfsClient(ipfsGateway);

            await ipfs.PubSub.SubscribeAsync(topic, async message => {
                var bytes = message.DataBytes;
                var item = JsonConvert.DeserializeObject<ISigned>(Encoding.UTF8.GetString(bytes));
                await output.AddAsync(item);
            }, CancellationToken.None);
        }
    }
}