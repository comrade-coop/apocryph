using System;
using System.Collections.Generic;
using System.Linq;
using System.Numerics;
using System.Runtime.CompilerServices;
using System.Text;
using System.Text.Json;
using System.Threading;
using System.Threading.Channels;
using System.Threading.Tasks;
using Apocryph.Core.Consensus;
using Apocryph.Core.Consensus.Blocks;
using Apocryph.Core.Consensus.Blocks.Command;
using Apocryph.Core.Consensus.Blocks.Messages;
using Apocryph.Core.Consensus.Communication;
using Apocryph.Core.Consensus.VirtualNodes;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Runtime.FunctionApp
{
    public class Proposer
    {
        private readonly Channel<Query<Block>> _channel;

        private Guid _chainId;
        private Guid _proposerAccount;
        private Node? _node;
        private Node[] _nodes;
        private IAsyncCollector<object>? _output;
        private Snowball<Block>? _snowball;
        private Node[]? _proposers;
        private Block? _lastBlock;
        private List<object>? _pendingCommands;
        private TaskCompletionSource<bool>? _pendingCommandsTaskCompletionSource;

        public Proposer()
        {
            _channel = Channel.CreateUnbounded<Query<Block>>();
        }

        [FunctionName(nameof(Proposer))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("chainId")] Guid chainId,
            [Perper("proposerAccount")] Guid proposerAccount,
            [PerperStream("assigner")] IAsyncEnumerable<(Node, Node[])> assigner,
            [Perper("lastBlock")] Block lastBlock,
            [Perper("pendingCommands")] List<object> pendingCommands,
            [PerperStream("gossips")] IAsyncEnumerable<Block> gossips,
            [PerperStream("queries")] IAsyncEnumerable<Query<Block>> queries,
            [PerperStream("output")] IAsyncCollector<object> output,
            CancellationToken cancellationToken)
        {
            _chainId = chainId;
            _proposerAccount = proposerAccount;
            _output = output;
            _lastBlock = lastBlock;
            _pendingCommands = pendingCommands;

            await Task.WhenAll(
                HandleAssigner(assigner),
                RunSnowball(context, cancellationToken),
                HandleIBC(gossips, cancellationToken),
                HandleQueries(queries, cancellationToken));
        }

        private async Task HandleAssigner(IAsyncEnumerable<(Node, Node[])> assigner)
        {
            await foreach (var (node, nodes) in assigner)
            {
                _node = node;
                _nodes = nodes;
            }
        }

        private async Task<Block> Propose(PerperStreamContext context)
        {
            var executor = new Executor(_node.Chain,
                async input => await context.CallWorkerAsync<(byte[]?, (string, object[])[], IDictionary<Guid, string[]>, IDictionary<Guid, string>)>("AgentWorker", new { input }, default));
            if (_pendingCommands!.Count == 0)
            {
                _pendingCommandsTaskCompletionSource = new TaskCompletionSource<bool>();
                // TODO: Possible race condition if TrySetResult happens before assigning a new completion source
                await _pendingCommandsTaskCompletionSource.Task;
                _pendingCommandsTaskCompletionSource = null;
            }

            var inputCommands = _pendingCommands.ToArray();
            _pendingCommands.Clear();

            if (_chainId == Guid.Empty)
            {
                inputCommands = inputCommands.Concat(new object[] {
                    new Invoke(_proposerAccount, (
                        "Apocryph.AgentZero.Messages.ClaimRewardMessage, Apocryph.AgentZero",
                        Encoding.UTF8.GetBytes("{}")))
                }).ToArray();
            }

            var (newState, newCommands, newCapabilities) = await executor.Execute(
                _lastBlock!.States, inputCommands, _lastBlock.Capabilities);
            // Include historical blocks as per protocol
            return new Block(_chainId, _proposerAccount, newState, inputCommands, newCommands, newCapabilities);
        }

        private async Task HandleIBC(IAsyncEnumerable<Block> ibc, CancellationToken cancellationToken)
        {
            var executor = new Executor(_node?.ToString()!, default!);

            await foreach (var block in ibc.WithCancellation(cancellationToken))
            {
                _pendingCommandsTaskCompletionSource?.TrySetResult(true);
                _pendingCommands!.AddRange(block.Commands.Where(x => executor.FilterCommand(x, _lastBlock!.Capabilities)));

                if (_chainId == Guid.Empty)
                {
                    _pendingCommands!.Add(new Invoke(_proposerAccount, (
                        typeof(SetChainBlockMessage).FullName!,
                        JsonSerializer.SerializeToUtf8Bytes(new SetChainBlockMessage
                        {
                            ChainId = block.ChainId,
                            BlockId = new byte[] { },
                            ProcessedCommands = new Dictionary<Guid, BigInteger>()
                            {
                                [block.ProposerAccount] = block.InputCommands.Length,
                            },
                            UsedTickets = new Dictionary<Guid, BigInteger>() { }, // TODO: Keep track of tickets
                            UnlockedTickets = new Dictionary<Guid, BigInteger>() { },
                        }))));
                }
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
                var proposerCount = _nodes.Length / 10; // TODO: Move constant to parameter
                _proposers = _nodes.Take(proposerCount).ToArray(); // TODO: Do a random walk based on last block hash

                var opinion = default(Block?);
                if (Array.IndexOf(_proposers, _node) != -1)
                {
                    opinion = await Propose(context);
                }

                _snowball = new Snowball<Block>(_node!, 100, 0.6, 3,
                                                SnowballSend, SnowballRespond, opinion);

                var committedProposal = await _snowball!.Run(_nodes, cancellationToken);

                await _output!.AddAsync(new Message<Block>(
                    committedProposal, MessageType.Committed), cancellationToken);
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