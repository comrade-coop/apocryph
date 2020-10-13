using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Core.Consensus;
using Apocryph.Core.Consensus.Blocks;
using Apocryph.Core.Consensus.Blocks.Command;
using Apocryph.Core.Consensus.Communication;
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
        private Snowball<Hash>? _snowball;
        private Dictionary<int, Hash> _pastBlocks = new Dictionary<int, Hash>();
        private int _round = 0;
        private Proposer? _proposer;
        private Func<Hash, Task<Block>>? _hashRegistryWorker;
        private IAsyncCollector<object>? _output;
        private Dictionary<Node, TaskCompletionSource<Query<Hash>>> _receiveCompletionSources = new Dictionary<Node, TaskCompletionSource<Query<Hash>>>();

        [FunctionName(nameof(ConsensusStream))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("proposerAccount")] Guid proposerAccount,
            [Perper("node")] Node node,
            [Perper("nodes")] List<Node> nodes,
            [Perper("chainData")] Chain chainData,
            [Perper("filter")] IAsyncEnumerable<Hash> filter,
            [Perper("validator")] IAsyncEnumerable<Message<Hash>> validator,
            [Perper("chain")] IAsyncEnumerable<Message<(Guid, Node?[])>> chain,
            [Perper("queries")] IAsyncEnumerable<Query<Hash>> queries,
            [Perper("hashRegistryWorker")] string hashRegistryWorker,
            [Perper("output")] IAsyncCollector<object> output,
            CancellationToken cancellationToken)
        {
            _output = output;
            _node = node;
            _nodes = nodes.ToArray();
            _hashRegistryWorker = hash => context.CallWorkerAsync<Block>(hashRegistryWorker, new { hash }, default);

            var executor = new Executor(_node!.ChainId,
                async (worker, input) => await context.CallWorkerAsync<(byte[]?, (string, object[])[], Dictionary<Guid, string[]>, Dictionary<Guid, string>)>(worker, new { input }, default));
            _proposer = new Proposer(executor, _node!.ChainId, chainData.GenesisBlock, new HashSet<Hash>(), new HashSet<ICommand>(), _node, proposerAccount);

            await TaskHelper.WhenAllOrFail(
                RunReports(context, cancellationToken),
                RunSnowball(context, cancellationToken),
                HandleChain(chain, cancellationToken),
                HandleValidator(validator, cancellationToken),
                HandleFilter(filter, cancellationToken),
                HandleQueries(queries, cancellationToken));
        }

        private async Task<Hash> Propose(PerperStreamContext context)
        {
            var proposal = await _proposer!.Propose();

            await _output!.AddAsync(proposal);

            return Hash.From(proposal);
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

        private async Task HandleValidator(IAsyncEnumerable<Message<Hash>> validator, CancellationToken cancellationToken)
        {
            await foreach (var message in validator.WithCancellation(cancellationToken))
            {
                if (message.Type == MessageType.Valid)
                {
                    var block = await _hashRegistryWorker!(message.Value);
                    _proposer!.AddConfirmedBlock(block!);
                }
            }
        }

        private async Task HandleFilter(IAsyncEnumerable<Hash> filter, CancellationToken cancellationToken)
        {
            await foreach (var hash in filter.WithCancellation(cancellationToken))
            {
                var block = await _hashRegistryWorker!(hash);
                _proposer!.AddConfirmedBlock(block!);
            }
        }

        private async Task HandleQueries(IAsyncEnumerable<Query<Hash>> queries, CancellationToken cancellationToken)
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
                            var reply = new Query<Hash>(_pastBlocks[query.Round], query.Receiver, query.Sender, query.Round, QueryVerb.Response);
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
                await Task.Delay(1500);

                var proposerCount = _nodes!.Length / 10 + 1; // TODO: Move constant to parameter
                _proposers = RandomWalk.Run(_proposer!.GetLastBlock()).Select(selected => _nodes[(int)(selected.Item1 % _nodes.Length)]).Take(proposerCount).ToArray();

                await _output!.AddAsync(new ConsensusReport(_node!, _proposers, _round, _proposer!.GetLastBlock()), cancellationToken);

                var opinion = default(Hash?);
                if (Array.IndexOf(_proposers, _node) != -1)
                {
                    opinion = await Propose(context);
                }

                await Task.Delay(500);

                var k = _nodes!.Length / 3 + 2; // TODO: Move constants to parameters
                var beta = (int)(Math.Log(_nodes!.Length) * 5) + 1; // TODO: Move constants to parameters

                _snowball = new Snowball<Hash>(_node!, k, 0.6, beta, _round,
                                                SnowballSend, SnowballRespond, opinion);

                var committedProposal = await _snowball!.Run(_nodes, cancellationToken);

                _pastBlocks[_round] = committedProposal;
                _round++;

                await _output!.AddAsync(new Message<Hash>(committedProposal, MessageType.Proposed), cancellationToken);
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

        private async Task<Query<Hash>> SnowballSend(Query<Hash> query, CancellationToken cancellationToken)
        {
            var taskCompletionSource = new TaskCompletionSource<Query<Hash>>();
            cancellationToken.Register(() => taskCompletionSource.TrySetCanceled());
            _receiveCompletionSources[query.Receiver] = taskCompletionSource;
            await _output!.AddAsync(query, cancellationToken);
            return await taskCompletionSource.Task;
        }

        private Query<Hash> SnowballRespond(Query<Hash> query, Hash? value)
        {
            var result = (value is null || IsNewBlockBetter(value.Value, query.Value) ? query.Value : value.Value);
            return new Query<Hash>(result, query.Receiver, query.Sender, query.Round, QueryVerb.Response);
        }

        private bool IsNewBlockBetter(Hash current, Hash suggested)
        {
            // TODO: change .Result to await
            var currentBlock = _hashRegistryWorker!(current).Result;
            var suggestedBlock = _hashRegistryWorker!(suggested).Result;

            if (currentBlock == null) return true;
            if (suggestedBlock == null) return false;

            var currentProposerOrder = Array.IndexOf(_proposers!, currentBlock.Proposer);
            var suggestedProposerOrder = Array.IndexOf(_proposers!, suggestedBlock.Proposer);

            return suggestedProposerOrder != -1 && suggestedProposerOrder < currentProposerOrder;
        }
    }
}