using System;
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
        private readonly HashSet<Hash> _finalizedBlocks = new HashSet<Hash>();
        private Node? _node;
        private Dictionary<Guid, Node?[]>? _nodes;
        private Committer? _committer;
        private IAsyncCollector<object>? _output;

        [FunctionName(nameof(IBCStream))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("node")] Node? node,
            [Perper("nodes")] Dictionary<Guid, Node?[]> nodes,
            [Perper("chain")] IAsyncEnumerable<Message<(Guid, Node?[])>> chain,
            [Perper("validator")] IAsyncEnumerable<Message<Hash>> validator,
            [Perper("gossips")] IAsyncEnumerable<Gossip<Hash>> gossips,
            [Perper("output")] IAsyncCollector<object> output,
            CancellationToken cancellationToken)
        {
            _output = output;
            _node = node;
            _nodes = nodes;
            _committer = new Committer();

            await TaskHelper.WhenAllOrFail(
                HandleChain(chain, cancellationToken),
                HandleValidator(validator, cancellationToken),
                HandleGossips(gossips, cancellationToken));
        }

        private async Task HandleValidator(IAsyncEnumerable<Message<Hash>> validator,
            CancellationToken cancellationToken)
        {
            await foreach (var message in validator)
            {
                var block = message.Value;
                var isValid = message.Type == MessageType.Valid;

                // Console.WriteLine("{0} sends gossip {1}", _node!, isValid);

                await _output!.AddAsync(new Gossip<Hash>(block, _node!,
                isValid ? GossipVerb.Confirm : GossipVerb.Reject), cancellationToken);
            }
        }

        private async Task HandleChain(IAsyncEnumerable<Message<(Guid, Node?[])>> chain, CancellationToken cancellationToken)
        {
            await foreach (var message in chain.WithCancellation(cancellationToken))
            {
                var (chainId, nodes) = message.Value;

                _nodes![chainId] = nodes;
            }
        }

        private async Task HandleGossips(IAsyncEnumerable<Gossip<Hash>> gossips,
            CancellationToken cancellationToken)
        {
            await foreach (var gossip in gossips.WithCancellation(cancellationToken))
            {
                if (!_nodes![gossip.Sender.ChainId].Contains(gossip.Sender) || _finalizedBlocks.Contains(gossip.Value))
                    continue;

                var nodes = _nodes[gossip.Sender.ChainId];

                _committer!.AddGossip(gossip);

                // if (_node is null) Console.WriteLine("got gossip {0} ({1})", string.Join(",", _committer!.GetConfirmations(gossip.Value, GossipVerb.Confirm, nodes)), string.Join(",", _committer!.GetConfirmations(gossip.Value, GossipVerb.Reject, nodes)));

                if (_committer!.IsGossipConfirmed(gossip.Value, GossipVerb.Reject, nodes))
                {
                    _finalizedBlocks.Add(gossip.Value);
                    await _output!.AddAsync(new Message<Hash>(gossip.Value, MessageType.Invalid), cancellationToken);
                }
                else if (_committer!.IsGossipConfirmed(gossip.Value, GossipVerb.Confirm, nodes))
                {
                    _finalizedBlocks.Add(gossip.Value);
                    await _output!.AddAsync(new Message<Hash>(gossip.Value, MessageType.Accepted), cancellationToken);
                }

                // Forward gossip
            }
        }
    }
}