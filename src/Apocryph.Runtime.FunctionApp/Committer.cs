using System.Collections.Generic;
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
    public class Committer
    {
        private readonly Dictionary<Block, bool> _validBlocks;
        private Block? _committedBlock;
        private readonly TaskCompletionSource<bool> _committedBlockValidTaskCompletionSource;

        private Node? _node;
        private IAsyncCollector<object>? _output;

        public Committer()
        {
            _validBlocks = new Dictionary<Block, bool>();
            _committedBlockValidTaskCompletionSource = new TaskCompletionSource<bool>();
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
                UpdateValidBlocks(validator, cancellationToken));
        }

        private async Task HandleProposals(IAsyncEnumerable<Message<Block>> proposer,
            CancellationToken cancellationToken)
        {
            await foreach (var proposerMessage in proposer)
            {
                _committedBlock = proposerMessage.Value;

                var committedBlockValid = _validBlocks.TryGetValue(_committedBlock, out var value)
                    ? value
                    : await _committedBlockValidTaskCompletionSource.Task;

                _committedBlock = null;

                await _output!.AddAsync(new Gossip<Block>(_committedBlock!, _node!,
                committedBlockValid ? GossipVerb.Confirm : GossipVerb.Reject), cancellationToken);
            }
        }

        private async Task UpdateValidBlocks(IAsyncEnumerable<Message<Block>> validator,
            CancellationToken cancellationToken)
        {
            await foreach (var message in validator.WithCancellation(cancellationToken))
            {
                _validBlocks[message.Value] = message.Type == MessageType.Valid;

                if (_committedBlock != null && _validBlocks.TryGetValue(_committedBlock, out var committedBlockValid))
                {
                    _committedBlockValidTaskCompletionSource.TrySetResult(committedBlockValid);
                }
            }
        }
    }
}