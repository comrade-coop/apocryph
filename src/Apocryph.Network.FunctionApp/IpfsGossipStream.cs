using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using System.Text.Json;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;
using Apocryph.Core.Consensus;
using Apocryph.Core.Consensus.Blocks;
using Apocryph.Core.Consensus.Communication;
using Apocryph.Core.Consensus.Serialization;
using Apocryph.Core.Consensus.VirtualNodes;
using Ipfs.Http;

namespace Apocryph.Runtime.FunctionApp
{
    public class IpfsGossipStream
    {
        static public readonly string PubsubChannel = "x-apocryph-gossip-0.0";

        private Dictionary<Node, Peer> nodeMappings = new Dictionary<Node, Peer>();
        private IAsyncCollector<Gossip<Hash>>? _output;

        [FunctionName(nameof(IpfsGossipStream))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("chain")] IAsyncEnumerable<(Node, Peer)> chain,
            [Perper("gossips")] IAsyncEnumerable<Gossip<Hash>> gossips,
            [Perper("output")] IAsyncCollector<Gossip<Hash>> output,
            CancellationToken cancellationToken)
        {
            _output = output;

            var ipfs = new IpfsClient();

            await TaskHelper.WhenAllOrFail(
                HandleChain(chain, cancellationToken),
                HandleGossips(ipfs, gossips, cancellationToken),
                RunSubscriber(ipfs, cancellationToken));
        }

        private async Task HandleChain(IAsyncEnumerable<(Node, Peer)> chain, CancellationToken cancellationToken)
        {
            await foreach (var (node, peer) in chain.WithCancellation(cancellationToken))
            {
                nodeMappings[node] = peer;
            }
        }

        private async Task HandleGossips(IpfsClient ipfs, IAsyncEnumerable<Gossip<Hash>> gossips, CancellationToken cancellationToken)
        {
            await foreach (var gossip in gossips.WithCancellation(cancellationToken))
            {
                var gossipBytes = JsonSerializer.Serialize(gossip, ApocryphSerializationOptions.JsonSerializerOptions);

                await ipfs.PubSub.PublishAsync(PubsubChannel, gossipBytes, cancellationToken);
            }
        }

        private async Task RunSubscriber(IpfsClient ipfs, CancellationToken cancellationToken)
        {
            await ipfs.PubSub.SubscribeAsync(PubsubChannel, async message =>
            {
                var gossip = JsonSerializer.Deserialize<Gossip<Hash>>(message.DataBytes, ApocryphSerializationOptions.JsonSerializerOptions);
                var peer = new Peer(message.Sender.Id.ToArray());

                if (nodeMappings.ContainsKey(gossip.Sender) && nodeMappings[gossip.Sender].Equals(peer))
                {
                    await _output!.AddAsync(gossip);
                }
            }, cancellationToken);
        }
    }
}