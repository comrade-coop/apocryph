using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Runtime.FunctionApp.Consensus.Core;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Runtime.FunctionApp.Consensus
{
    public class Committer
    {
        private readonly Dictionary<Block, bool> _validBlocks;
        private Block? _committedBlock;
        private readonly TaskCompletionSource<bool> _committedBlockValidTaskCompletionSource;
        private bool _accepted;

        private Node? _node;
        private IAsyncCollector<object>? _output;

        public Committer()
        {
            _validBlocks = new Dictionary<Block, bool>();
            _committedBlockValidTaskCompletionSource = new TaskCompletionSource<bool>();
            _accepted = false;
        }

        [FunctionName(nameof(Committer))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("node")] Node node,
            [Perper("nodes")] Node[] nodes,
            [PerperStream("gossips")] IAsyncEnumerable<Gossip<Block>> gossips,
            [PerperStream("proposer")] IAsyncEnumerable<Message<Block>> proposer,
            [PerperStream("validator")] IAsyncEnumerable<Message<Block>> validator,
            [PerperStream("output")] IAsyncCollector<object> output,
            CancellationToken cancellationToken)
        {
            _node = node;
            _output = output;

            await Task.WhenAll(
                HandleProposals(proposer, cancellationToken),
                HandleGossips(gossips, nodes, cancellationToken),
                UpdateValidBlocks(validator, cancellationToken));
        }

        private async Task HandleProposals(IAsyncEnumerable<Message<Block>> proposer,
            CancellationToken cancellationToken)
        {
            var commitQuery = await proposer.FirstAsync(cancellationToken);
            _committedBlock = commitQuery.Value;

            var committedBlockValid = _validBlocks.TryGetValue(_committedBlock, out var value)
                ? value
                : await _committedBlockValidTaskCompletionSource.Task;
            await _output!.AddAsync(new Gossip<Block>(_committedBlock, new[] { _node! },
                committedBlockValid ? GossipVerb.Confirm : GossipVerb.Reject), cancellationToken);
        }

        private async Task HandleGossips(IAsyncEnumerable<Gossip<Block>> gossips,
            Node[] nodes,
            CancellationToken cancellationToken)
        {
            await foreach (var gossip in gossips.WithCancellation(cancellationToken))
            {
                if (gossip.Verb == GossipVerb.Confirm &&
                    gossip.Signers.Contains(_node) &&
                    3 * gossip.Signers.Length > 2 * nodes.Length)
                {
                    _accepted = true;
                    await _output!.AddAsync(new Message<Block>(gossip.Value, MessageType.Accepted), cancellationToken);
                    break;
                }

                if (gossip.Signers.Contains(_node)) continue;

                var committedBlockValid = _validBlocks.TryGetValue(_committedBlock!, out var value)
                    ? value
                    : await _committedBlockValidTaskCompletionSource.Task;
                if (gossip.Value.Equals(_committedBlock!) && committedBlockValid)
                {
                    var signers = gossip.Verb == GossipVerb.Confirm
                        ? gossip.Signers.Append(_node!).ToArray()
                        : new[] { _node! };
                    await _output!.AddAsync(new Gossip<Block>(gossip.Value, signers, GossipVerb.Confirm),
                        cancellationToken);
                }
                else
                {
                    var signers = gossip.Verb == GossipVerb.Reject
                        ? gossip.Signers.Append(_node!).ToArray()
                        : new[] { _node! };
                    await _output!.AddAsync(new Gossip<Block>(gossip.Value, signers, GossipVerb.Reject),
                        cancellationToken);
                }
            }
        }

        private async Task UpdateValidBlocks(IAsyncEnumerable<Message<Block>> validator,
            CancellationToken cancellationToken)
        {
            await foreach (var message in validator.WithCancellation(cancellationToken))
            {
                if (_accepted) break;

                _validBlocks[message.Value] = message.Type == MessageType.Valid;

                if (_committedBlock != null && _validBlocks.TryGetValue(_committedBlock, out var committedBlockValid))
                {
                    _committedBlockValidTaskCompletionSource.TrySetResult(committedBlockValid);
                }
            }
        }
    }
}