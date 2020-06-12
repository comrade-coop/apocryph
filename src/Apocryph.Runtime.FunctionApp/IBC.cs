using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Core.Consensus.Blocks;
using Apocryph.Core.Consensus.Communication;
using Apocryph.Core.Consensus.VirtualNodes;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Runtime.FunctionApp
{
    public class IBC
    {
        private readonly Dictionary<Block, HashSet<Node>> _gossipConfirmations = new Dictionary<Block, HashSet<Node>>();
        private readonly HashSet<Block> _acceptedBlocks = new HashSet<Block>();
        private Message<Block>? _committingMessage;
        private Node? _node;
        private Node[]? _nodes;
        private IAsyncCollector<object>? _output;

        [FunctionName(nameof(IBC))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("node")] Node node,
            [Perper("nodes")] Node[] nodes,
            [PerperStream("validator")] IAsyncEnumerable<Message<Block>> validator,
            [PerperStream("gossips")] IAsyncEnumerable<Gossip<Block>> gossips,
            [PerperStream("output")] IAsyncCollector<object> output,
            CancellationToken cancellationToken)
        {
            _output = output;
            _node = node;
            _nodes = nodes;

            await Task.WhenAll(
                HandleValidator(validator, cancellationToken),
                HandleGossips(gossips, cancellationToken));
        }

        private async Task HandleValidator(IAsyncEnumerable<Message<Block>> validator,
            CancellationToken cancellationToken)
        {
            await foreach (var message in validator)
            {
                var block = message.Value;
                var isValid = message.Type == MessageType.Valid;

                _committingMessage = message;

                await _output!.AddAsync(new Gossip<Block>(block, _node!,
                isValid ? GossipVerb.Confirm : GossipVerb.Reject), cancellationToken);
            }
        }

        private async Task HandleGossips(IAsyncEnumerable<Gossip<Block>> gossips,
            CancellationToken cancellationToken)
        {
            await foreach (var gossip in gossips.WithCancellation(cancellationToken))
            {
                if (!_nodes!.Contains(gossip.Sender) || _acceptedBlocks.Contains(gossip.Value))
                    continue;

                if (gossip.Verb == GossipVerb.Confirm) // TODO: Count rejections
                {
                    if (!_gossipConfirmations.ContainsKey(gossip.Value))
                    {
                        _gossipConfirmations[gossip.Value] = new HashSet<Node>();
                    }
                    var confirmations = _gossipConfirmations[gossip.Value];

                    confirmations.Add(gossip.Sender);

                    if (3 * confirmations.Count > 2 * _nodes!.Length)
                    {
                        _acceptedBlocks.Add(gossip.Value);
                        await _output!.AddAsync(new Message<Block>(gossip.Value, MessageType.Accepted), cancellationToken);
                    }
                }
                // Forward gossip
            }
        }
    }
}