using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Core.Consensus;
using Apocryph.Core.Consensus.Blocks;
using Apocryph.Core.Consensus.Communication;
using Apocryph.Core.Consensus.VirtualNodes;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Runtime.FunctionApp
{
    public class IBCStream
    {
        private readonly HashSet<Block> _finalizedBlocks = new HashSet<Block>();
        private Node? _node;
        private Dictionary<byte[], Node?[]>? _nodes;
        private Committer? _committer;
        private IAsyncCollector<object>? _output;

        [FunctionName(nameof(IBCStream))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("node")] Node node,
            [Perper("nodes")] Dictionary<byte[], Node?[]> nodes,
            [PerperStream("queries")] IAsyncEnumerable<Message<(byte[], Node?[])>> chain,
            [PerperStream("validator")] IAsyncEnumerable<Message<Block>> validator,
            [PerperStream("gossips")] IAsyncEnumerable<Gossip<Block>> gossips,
            [PerperStream("output")] IAsyncCollector<object> output,
            CancellationToken cancellationToken)
        {
            _output = output;
            _node = node;
            _nodes = nodes;
            _committer = new Committer();

            await Task.WhenAll(
                HandleChain(chain, cancellationToken),
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

                await _output!.AddAsync(new Gossip<Block>(block, _node!,
                isValid ? GossipVerb.Confirm : GossipVerb.Reject), cancellationToken);
            }
        }

        private async Task HandleChain(IAsyncEnumerable<Message<(byte[], Node?[])>> chain, CancellationToken cancellationToken)
        {
            await foreach (var message in chain.WithCancellation(cancellationToken))
            {
                var (chainId, nodes) = message.Value;

                _nodes![chainId] = nodes;
            }
        }

        private async Task HandleGossips(IAsyncEnumerable<Gossip<Block>> gossips,
            CancellationToken cancellationToken)
        {
            await foreach (var gossip in gossips.WithCancellation(cancellationToken))
            {
                if (!_nodes![gossip.Sender.ChainId].Contains(gossip.Sender) || _finalizedBlocks.Contains(gossip.Value))
                    continue;

                var nodes = _nodes[gossip.Sender.ChainId];

                _committer!.AddGossip(gossip);

                if (_committer!.IsGossipConfirmed(gossip.Value, GossipVerb.Reject, nodes))
                {
                    _finalizedBlocks.Add(gossip.Value);
                    await _output!.AddAsync(new Message<Block>(gossip.Value, MessageType.Invalid), cancellationToken);
                }
                else if (_committer!.IsGossipConfirmed(gossip.Value, GossipVerb.Confirm, nodes))
                {
                    _finalizedBlocks.Add(gossip.Value);
                    await _output!.AddAsync(new Message<Block>(gossip.Value, MessageType.Accepted), cancellationToken);
                }

                // Forward gossip
            }
        }
    }
}