using System;
using System.Collections.Generic;
using System.Linq;
using System.Runtime.CompilerServices;
using System.Threading;
using System.Threading.Channels;
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
    public class ConsensusStream
    {
        private readonly Channel<Query<Block>> _channel;

        private Node? _node;
        private Node[]? _nodes;
        private Node[]? _proposers;
        private Snowball<Block>? _snowball;
        private Proposer? _proposer;
        private IAsyncCollector<object>? _output;

        public ConsensusStream()
        {
            _channel = Channel.CreateUnbounded<Query<Block>>();
        }

        [FunctionName(nameof(ConsensusStream))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("node")] Node node,
            [Perper("nodes")] Node[] nodes,
            [Perper("proposerAccount")] Guid proposerAccount,
            [Perper("lastBlock")] Block lastBlock,
            [Perper("pendingCommands")] HashSet<object> pendingCommands,
            [PerperStream("filter")] IAsyncEnumerable<Block> filter,
            [PerperStream("chain")] IAsyncEnumerable<Message<(byte[], Node?[])>> chain,
            [PerperStream("queries")] IAsyncEnumerable<Query<Block>> queries,
            [PerperStream("output")] IAsyncCollector<object> output,
            CancellationToken cancellationToken)
        {
            _output = output;
            _node = node;
            _nodes = nodes;

            var executor = new Executor(_node!.ChainId,
                async input => await context.CallWorkerAsync<(byte[]?, (string, object[])[], IDictionary<Guid, string[]>, IDictionary<Guid, string>)>("AgentWorker", new { input }, default));
            _proposer = new Proposer(executor, _node!.ChainId, lastBlock, pendingCommands, proposerAccount);

            await Task.WhenAll(
                RunSnowball(context, cancellationToken),
                HandleChain(chain, cancellationToken),
                HandleFilter(filter, cancellationToken),
                HandleQueries(queries, cancellationToken));
        }

        private Task<Block> Propose(PerperStreamContext context)
        {
            return _proposer!.Propose();
        }

        private async Task HandleChain(IAsyncEnumerable<Message<(byte[], Node?[])>> chain, CancellationToken cancellationToken)
        {
            await foreach (var message in chain.WithCancellation(cancellationToken))
            {
                var (chainId, nodes) = message.Value;

                if (_node!.ChainId.SequenceEqual(chainId))
                {
                    _nodes = nodes as Node[]; // TODO: Remove cast
                }
            }
        }

        private async Task HandleFilter(IAsyncEnumerable<Block> filter, CancellationToken cancellationToken)
        {
            await foreach (var block in filter.WithCancellation(cancellationToken))
            {
                _proposer!.AddConfirmedBlock(block);
            }
        }

        private async Task HandleQueries(IAsyncEnumerable<Query<Block>> queries, CancellationToken cancellationToken)
        {
            await foreach (var query in queries.WithCancellation(cancellationToken))
            {
                if (_snowball is null) continue;

                if (query.Receiver == _node && query.Verb == QueryVerb.Response)
                {
                    await _channel.Writer.WriteAsync(query, cancellationToken);
                }
                else if (query.Receiver == _node && query.Verb == QueryVerb.Request)
                {
                    await _output!.AddAsync(_snowball!.Query(query), cancellationToken);
                }
            }
        }

        private async Task RunSnowball(PerperStreamContext context, CancellationToken cancellationToken)
        {
            while (true)
            {
                var proposerCount = _nodes!.Length / 10; // TODO: Move constant to parameter
                _proposers = _nodes.Take(proposerCount).ToArray(); // TODO: Do a random walk based on last block hash

                var opinion = default(Block?);
                if (Array.IndexOf(_proposers, _node) != -1)
                {
                    opinion = await Propose(context);
                }

                _snowball = new Snowball<Block>(_node!, 100, 0.6, 3,
                                                SnowballSend, SnowballRespond, opinion);

                var committedProposal = await _snowball!.Run(_nodes, cancellationToken);

                await _output!.AddAsync(new Message<Block>(committedProposal, MessageType.Proposed), cancellationToken);
            }
        }

        private async IAsyncEnumerable<Query<Block>> SnowballSend(Query<Block>[] queries, [EnumeratorCancellation] CancellationToken cancellationToken)
        {
            await Task.WhenAll(queries.Select(q => _output!.AddAsync(q, cancellationToken)));
            await foreach (var answer in _channel.Reader.ReadAllAsync(cancellationToken).Take(queries.Length).WithCancellation(cancellationToken))
            {
                yield return answer;
            }
        }

        private Query<Block> SnowballRespond(Query<Block> query, Block? value, Block? opinion)
        {
            var result = opinion ??
                         (value is null || IsNewBlockBetter(value, query.Value) ? query.Value : value);
            return new Query<Block>(result, query.Receiver, query.Sender, QueryVerb.Response);
        }

        private bool IsNewBlockBetter(Block current, Block suggested)
        {
            var currentProposerOrder = Array.IndexOf(_proposers!, current.Proposer);
            var suggestedProposerOrder = Array.IndexOf(_proposers!, suggested.Proposer);

            return suggestedProposerOrder != -1 && suggestedProposerOrder < currentProposerOrder;
        }
    }
}