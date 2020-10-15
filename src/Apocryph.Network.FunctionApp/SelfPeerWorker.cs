using System.Threading.Tasks;
using Apocryph.Core.Consensus.VirtualNodes;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;
using Ipfs.Http;

namespace Apocryph.Runtime.FunctionApp
{
    public class SelfPeerWorker
    {
        [FunctionName(nameof(SelfPeerWorker))]
        [return: Perper("$return")]
        public async Task<Peer> Run([PerperWorkerTrigger] PerperWorkerContext context)
        {
            var ipfsClient = new IpfsClient();

            return new Peer((await ipfsClient.IdAsync()).Id.ToArray());
        }
    }
}