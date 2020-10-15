using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using System.Text.Json;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;
using Apocryph.Core.Consensus;
using Apocryph.Core.Consensus.Serialization;
using Apocryph.Core.Consensus.VirtualNodes;
using Ipfs.Http;

namespace Apocryph.Runtime.FunctionApp
{
    public class IpfsSlotGossipStream
    {
        public class SlotGossip
        {
            public Guid ChainId { get; set; }
            public byte[] Proof { get; set; }

            public SlotGossip(Guid chainId, byte[] proof)
            {
                ChainId = chainId;
                Proof = proof;
            }
        }

        static public readonly string PubsubChannel = "x-apocryph-slot-gossip-0.0";

        private IAsyncCollector<SlotClaim>? _output;

        [FunctionName(nameof(IpfsSlotGossipStream))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("self")] Peer self,
            [Perper("claims")] IAsyncEnumerable<SlotClaim> claims,
            [Perper("output")] IAsyncCollector<SlotClaim> output,
            CancellationToken cancellationToken)
        {
            _output = output;

            var ipfs = new IpfsClient();

            await TaskHelper.WhenAllOrFail(
                HandleGossips(ipfs, claims, cancellationToken),
                RunSubscriber(ipfs, cancellationToken));
        }

        private async Task HandleGossips(IpfsClient ipfs, IAsyncEnumerable<SlotClaim> claims, CancellationToken cancellationToken)
        {
            await foreach (var claim in claims.WithCancellation(cancellationToken))
            {
                var gossip = new SlotGossip(claim.ChainId, claim.Proof);
                var gossipBytes = JsonSerializer.Serialize(gossip, ApocryphSerializationOptions.JsonSerializerOptions);

                await ipfs.PubSub.PublishAsync(PubsubChannel, gossipBytes, cancellationToken);
            }
        }

        private async Task RunSubscriber(IpfsClient ipfs, CancellationToken cancellationToken)
        {
            await ipfs.PubSub.SubscribeAsync(PubsubChannel, async message =>
            {
                var gossip = JsonSerializer.Deserialize<SlotGossip>(message.DataBytes, ApocryphSerializationOptions.JsonSerializerOptions);
                var peer = new Peer(message.Sender.Id.ToArray());
                var claim = new SlotClaim(gossip.ChainId, peer, gossip.Proof);

                await _output!.AddAsync(claim);
            }, cancellationToken);
        }
    }
}