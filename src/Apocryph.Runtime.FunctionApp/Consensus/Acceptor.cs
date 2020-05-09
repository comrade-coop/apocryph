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
    public class Acceptor
    {
        private Dictionary<Block, bool> _validBlocks;
        private Block? _acceptedBlock;
        private bool _committed;

        private Node? _node;
        private IAsyncCollector<Gossip<Block>>? _output;

        public Acceptor()
        {
            _validBlocks = new Dictionary<Block, bool>();
            _committed = false;
        }

        [FunctionName(nameof(Acceptor))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("node")] Node node,
            [Perper("nodes")] Node[] nodes,
            [PerperStream("gossips")] IAsyncEnumerable<Gossip<Block>> gossips,
            [PerperStream("proposer")] IAsyncEnumerable<Query<Block>> proposer,
            [PerperStream("validator")] IAsyncEnumerable<Message<Block>> validator,
            [PerperStream("output")] IAsyncCollector<Gossip<Block>> output,
            CancellationToken cancellationToken)
        {
            _node = node;
            _output = output;

            await Task.WhenAll(
                HandleProposals(proposer, cancellationToken),
                HandleGossips(gossips, nodes, cancellationToken),
                UpdateValidBlocks(validator, cancellationToken));
        }

        private async Task HandleProposals(IAsyncEnumerable<Query<Block>> proposer,
            CancellationToken cancellationToken)
        {
            var acceptQuery = await proposer.FirstAsync(
                query => query.Sender == _node && query.Receiver == _node && query.Verb == QueryVerb.Accept,
                cancellationToken);
            _acceptedBlock = acceptQuery.Value;
            await _output!.AddAsync(new Gossip<Block>(_acceptedBlock, new[] {_node!}, GossipVerb.Confirm),
                cancellationToken);
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
                    _committed = true;
                    await _output!.AddAsync(new Gossip<Block>(gossip.Value, gossip.Signers, GossipVerb.Commit),
                        cancellationToken);
                    break;
                }

                if (gossip.Signers.Contains(_node)) continue;

                if (gossip.Value.Equals(_acceptedBlock!) && _validBlocks.GetValueOrDefault(_acceptedBlock!))
                {
                    var signers = gossip.Verb == GossipVerb.Confirm
                        ? gossip.Signers.Append(_node!).ToArray()
                        : new[] {_node!};
                    await _output!.AddAsync(new Gossip<Block>(gossip.Value, signers, GossipVerb.Confirm),
                        cancellationToken);
                }
                else
                {
                    var signers = gossip.Verb == GossipVerb.Reject
                        ? gossip.Signers.Append(_node!).ToArray()
                        : new[] {_node!};
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
                if (_committed) break;

                _validBlocks[message.Value] = message.Type == MessageType.Valid;
            }
        }
    }
}