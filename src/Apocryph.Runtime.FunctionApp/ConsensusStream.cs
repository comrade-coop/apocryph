using System;
using System.Collections.Generic;
using System.Linq;
using System.Text.Json;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Core.Consensus;
using Apocryph.Core.Consensus.Blocks;
using Apocryph.Core.Consensus.Blocks.Command;
using Apocryph.Core.Consensus.Communication;
using Apocryph.Core.Consensus.Serialization;
using Apocryph.Core.Consensus.VirtualNodes;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Runtime.FunctionApp
{
    public class ConsensusStream
    {
        private Node? _node;
        private Node?[]? _nodes;
        private Node?[]? _proposers;
        private Snowball<Block>? _snowball;
        private Dictionary<int, Block> _pastBlocks = new Dictionary<int, Block>();
        private int _round = 0;
        private Proposer? _proposer;
        private IAsyncCollector<object>? _output;
        private Dictionary<Node, TaskCompletionSource<Query<Block>>> _receiveCompletionSources = new Dictionary<Node, TaskCompletionSource<Query<Block>>>();

        [FunctionName(nameof(ConsensusStream))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("proposerAccount")] Guid proposerAccount,
            [Perper("node")] Node node,
            [Perper("nodes")] List<Node> nodes,
            [Perper("chainData")] Chain chainData,
            [Perper("filter")] IAsyncEnumerable<Block> filter,
            [Perper("validator")] IAsyncEnumerable<Message<Block>> validator,
            [Perper("chain")] IAsyncEnumerable<Message<(Guid, Node?[])>> chain,
            [Perper("queries")] IAsyncEnumerable<Query<Block>> queries,
            [Perper("output")] IAsyncCollector<object> output,
            CancellationToken cancellationToken)
        {
            _output = output;
            _node = node;
            _nodes = nodes.ToArray();

            var executor = new Executor(_node!.ChainId,
                async (worker, input) => await context.CallWorkerAsync<(byte[]?, (string, object[])[], Dictionary<Guid, string[]>, Dictionary<Guid, string>)>(worker, new { input }, default));
            _proposer = new Proposer(executor, _node!.ChainId, chainData.GenesisBlock, new HashSet<Block>(), new HashSet<ICommand>(), _node, proposerAccount);

            await TaskHelper.WhenAllOrFail(
                RunReports(context, cancellationToken),
                RunSnowball(context, cancellationToken),
                HandleChain(chain, cancellationToken),
                HandleValidator(validator, cancellationToken),
                HandleFilter(filter, cancellationToken),
                HandleQueries(queries, cancellationToken));
        }

        private Task<Block> Propose(PerperStreamContext context)
        {
            return _proposer!.Propose();
        }

        private async Task HandleChain(IAsyncEnumerable<Message<(Guid, Node?[])>> chain, CancellationToken cancellationToken)
        {
            await foreach (var message in chain.WithCancellation(cancellationToken))
            {
                var (chainId, nodes) = message.Value;

                if (_node!.ChainId == chainId)
                {
                    _nodes = nodes;
                }
            }
        }

        private async Task HandleValidator(IAsyncEnumerable<Message<Block>> validator, CancellationToken cancellationToken)
        {
            await foreach (var message in validator.WithCancellation(cancellationToken))
            {
                if (message.Type == MessageType.Valid)
                {
                    _proposer!.AddConfirmedBlock(message.Value);
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
                if (_snowball is null || !query.Receiver.Equals(_node)) continue;

                if (query.Verb == QueryVerb.Response)
                {
                    _receiveCompletionSources[query.Sender].TrySetResult(query);
                }
                else if (query.Verb == QueryVerb.Request)
                {
                    if (_round != query.Round)
                    {
                        if (_pastBlocks.ContainsKey(query.Round))
                        {
                            var reply = new Query<Block>(_pastBlocks[query.Round], query.Receiver, query.Sender, query.Round, QueryVerb.Response);
                            await _output!.AddAsync(reply, cancellationToken);
                        }
                    }
                    else
                    {
                        await _output!.AddAsync(_snowball!.Query(query), cancellationToken);
                    }
                }
            }
        }

        private async Task RunSnowball(PerperStreamContext context, CancellationToken cancellationToken)
        {
            await Task.Delay(4000);
            while (true)
            {
                await Task.Delay(2000);

                var proposerCount = _nodes!.Length / 10 + 1; // TODO: Move constant to parameter
                var serializedBlock = JsonSerializer.SerializeToUtf8Bytes(_proposer!.GetLastBlock(), ApocryphSerializationOptions.JsonSerializerOptions);
                _proposers = RandomWalk.Run(serializedBlock).Select(selected => _nodes[(int)(selected.Item1 % _nodes.Length)]).Take(proposerCount).ToArray();

                await _output!.AddAsync(new ConsensusReport(_node!, _proposers, _round, _proposer!.GetLastBlock()), cancellationToken);

                var opinion = default(Block?);
                if (Array.IndexOf(_proposers, _node) != -1)
                {
                    opinion = await Propose(context);
                }

                var k = _nodes!.Length / 3 + 2; // TODO: Move constants to parameters
                var beta = _nodes!.Length * 2; // TODO: Move constants to parameters

                // Console.WriteLine("{0} starts new round! {1} {2}", _node, _round, string.Join(",", (IEnumerable<Node?>)_proposers));
                _snowball = new Snowball<Block>(_node!, k, 0.6, beta, _round,
                                                SnowballSend, SnowballRespond, opinion);

                var committedProposal = await _snowball!.Run(_nodes, cancellationToken);

                _pastBlocks[_round] = committedProposal;
                _round++;

                await _output!.AddAsync(new Message<Block>(committedProposal, MessageType.Proposed), cancellationToken);
            }
        }

        private async Task RunReports(PerperStreamContext context, CancellationToken cancellationToken)
        {
            while (!cancellationToken.IsCancellationRequested)
            {
                await Task.Delay(1000, cancellationToken);

                if (_snowball != null)
                {
                    await _output!.AddAsync(new SnowballReport(_node!, _snowball.GetConfirmedValues()), cancellationToken);
                }
            }
        }

        private async Task<Query<Block>> SnowballSend(Query<Block> query, CancellationToken cancellationToken)
        {
            var taskCompletionSource = new TaskCompletionSource<Query<Block>>();
            cancellationToken.Register(() => taskCompletionSource.TrySetCanceled());
            _receiveCompletionSources[query.Receiver] = taskCompletionSource;
            await _output!.AddAsync(query, cancellationToken);
            return await taskCompletionSource.Task;
        }

        private Query<Block> SnowballRespond(Query<Block> query, Block? value)
        {
            var result = (value is null || IsNewBlockBetter(value, query.Value) ? query.Value : value);
            return new Query<Block>(result, query.Receiver, query.Sender, query.Round, QueryVerb.Response);
        }

        private bool IsNewBlockBetter(Block current, Block suggested)
        {
            var currentProposerOrder = Array.IndexOf(_proposers!, current.Proposer);
            var suggestedProposerOrder = Array.IndexOf(_proposers!, suggested.Proposer);

            return suggestedProposerOrder != -1 && suggestedProposerOrder < currentProposerOrder;
        }
    }
}