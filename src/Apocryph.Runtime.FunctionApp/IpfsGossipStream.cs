using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using System.Text.Json;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;
using Apocryph.Core.Consensus.Communication;
using Apocryph.Core.Consensus.Blocks;
using Apocryph.Core.Consensus.Serialization;
using Ipfs.Http;

namespace Apocryph.Runtime.FunctionApp
{
    public class IpfsGossipStream
    {
        // TODO: Rework to a custom gossipsub implementation on top of IPFS p2p.
        static public readonly string PubsubChannel = "x-apocryph-gossip-0.0";

        private IAsyncCollector<Gossip<Hash>>? _output ;

        [FunctionName(nameof(IpfsGossipStream))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("queries")] IAsyncEnumerable<Gossip<Hash>> gossips,
            [Perper("output")] IAsyncCollector<Gossip<Hash>> output,
            CancellationToken cancellationToken)
        {
            _output = output;
            var ipfs = new IpfsClient();
            await TaskHelper.WhenAllOrFail(
                HandleGossips(ipfs, gossips, cancellationToken),
                RunSubscriber(ipfs, cancellationToken));
        }

        private async Task HandleGossips(IpfsClient ipfs, IAsyncEnumerable<Gossip<Hash>> gossips, CancellationToken cancellationToken)
        {
            await foreach (var gossip in gossips.WithCancellation(cancellationToken))
            {
                var gossipBytes = JsonSerializer.Serialize(gossip, ApocryphSerializationOptions.JsonSerializerOptions);

                await ipfs.PubSub.PublishAsync(PubsubChannel, gossipBytes, CancellationToken.None);
            }
        }

        private async Task RunSubscriber(IpfsClient ipfs, CancellationToken cancellationToken)
        {
            await ipfs.PubSub.SubscribeAsync(PubsubChannel, async message =>
            {
                var gossip = JsonSerializer.Deserialize<Gossip<Hash>>(message.DataBytes, ApocryphSerializationOptions.JsonSerializerOptions);

                await _output!.AddAsync(gossip);
            }, CancellationToken.None);
        }
    }
}