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
    public class Acceptor
    {
        private readonly Dictionary<Block, HashSet<Node>> _gossipConfirmations = new Dictionary<Block, HashSet<Node>>();
        private readonly Dictionary<Block, bool> _validatedBlocks = new Dictionary<Block, bool>();
        private readonly HashSet<Block> _acceptedBlocks = new HashSet<Block>();

        private IAsyncCollector<Block>? _output;
        private Node _node;
        private Node[] _nodes;

        [FunctionName(nameof(Acceptor))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("assigner")] IAsyncEnumerable<(Node, Node[])> assigner,
            [PerperStream("gossips")] IAsyncEnumerable<Gossip<Block>> gossips,
            [PerperStream("output")] IAsyncCollector<Block> output,
            CancellationToken cancellationToken)
        {
            _output = output;

            await Task.WhenAll(
                HandleAssigner(assigner),
                HandleGossips(gossips, cancellationToken));
        }

        private async Task HandleAssigner(IAsyncEnumerable<(Node, Node[])> assigner)
        {
            await foreach (var (node, nodes) in assigner)
            {
                _node = node;
                _nodes = nodes;
            }
        }

        private async Task HandleGossips(IAsyncEnumerable<Gossip<Block>> gossips,
            CancellationToken cancellationToken)
        {
            await foreach (var gossip in gossips.WithCancellation(cancellationToken))
            {
                if (!_nodes.Contains(gossip.Sender) || _acceptedBlocks.Contains(gossip.Value))
                    continue;

                if (gossip.Verb == GossipVerb.Confirm)
                {
                    GetGossipConfirmations(gossip.Value).Add(gossip.Sender);

                    if (3 * GetGossipConfirmations(gossip.Value).Count > 2 * _nodes.Length)
                    {
                        _acceptedBlocks.Add(gossip.Value);
                        // TODO: Check block validity and forward gossip
                        await _output!.AddAsync(gossip.Value, cancellationToken);
                    }
                }

                // if (gossip.Verb == GossipVerb.IdentityChanged)
                // {
                //     GetGossipConfirmations(gossip.Value).Remove(gossip.Sender);
                // }
            }
        }

        private HashSet<Node> GetGossipConfirmations(Block block)
        {
            if (!_gossipConfirmations.ContainsKey(block))
            {
                _gossipConfirmations[block] = new HashSet<Node>();
            }
            return _gossipConfirmations[block];
        }
    }
}