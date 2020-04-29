using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Agent;
using Apocryph.Runtime.FunctionApp.Ipfs;
using Ipfs;
using Ipfs.Http;
using Microsoft.Azure.WebJobs;
using Newtonsoft.Json;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Runtime.FunctionApp.Commuication
{
    public static class SubscriptionCommandOutput
    {
        [FunctionName(nameof(SubscriptionCommandOutput))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("otherId")] string otherId,
            [PerperStream("publicationsStream")] IAsyncEnumerable<PublicationCommand> publicationsStream,
            [PerperStream("outputStream")] IAsyncCollector<(string, object)> outputStream)
        {
            await publicationsStream.ForEachAsync(async publication =>
            {
                await outputStream.AddAsync((otherId, publication.Payload));
            }, CancellationToken.None);
        }
    }
}