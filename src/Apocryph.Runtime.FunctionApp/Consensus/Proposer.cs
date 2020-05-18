using System;
using System.Collections.Generic;
using System.Runtime.CompilerServices;
using System.Linq;
using System.Threading;
using System.Threading.Channels;
using System.Threading.Tasks;
using Apocryph.Agent.Command;
using Apocryph.Agent.Core;
using Apocryph.Agent.Worker;
using Apocryph.Runtime.FunctionApp.Consensus.Core;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Runtime.FunctionApp.Consensus
{
    public class Proposer
    {
        private readonly Channel<Query<Block>> _channel;

        private Node? _node;
        private IAsyncCollector<object>? _output;
        private Snowball<Block>? _snowball;
        private Node[]? _proposers;
        private Block? _lastBlock;

        public Proposer()
        {
            _channel = Channel.CreateUnbounded<Query<Block>>();
        }

        [FunctionName(nameof(Proposer))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("node")] Node node,
            [Perper("nodes")] Node[] nodes,
            [Perper("lastBlock")] Block lastBlock,
            [PerperStream("queries")] IAsyncEnumerable<Query<Block>> queries,
            [PerperStream("output")] IAsyncCollector<object> output,
            CancellationToken cancellationToken)
        {
            _node = node;
            _output = output;
            _lastBlock = lastBlock;

            await Task.WhenAll(
                RunSnowball(context, nodes, cancellationToken),
                HandleQueries(queries, cancellationToken));
        }

        private async Task<Block> Propose(PerperStreamContext context)
        {
            var executor = new Executor(_node?.ToString()!,
                async input => await context.CallWorkerAsync<WorkerOutput>("AgentWorker", new {input}, default));
            var command = _lastBlock.Commands.FirstOrDefault(o => o is Invoke || o is Publish || o is Remind);
            if (command != null)
            {
                var (newState, newCommands, newCapabilities) = await executor.Execute(
                    _lastBlock.State, command, _lastBlock.Capabilities);
                return new Block(newState, newCommands, newCapabilities);
            }
            // Include historical blocks as per protocol
            return default!; //What to return?
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

        private async Task RunSnowball(PerperStreamContext context, Node[] nodes, CancellationToken cancellationToken)
        {
            while (true)
            {
                var proposerCount = nodes.Length / 10; // TODO: Move constant to parameter
                _proposers = nodes.Take(proposerCount).ToArray(); // TODO: Do a random walk based on last block hash

                var opinion = default(Block?);
                if (Array.IndexOf(_proposers, _node) != -1)
                {
                    opinion = await Propose(context);
                }

                _snowball = new Snowball<Block>(_node!, 100, 0.6, 3,
                                                SnowballSend, SnowballRespond, opinion);

                var committedProposal = await _snowball!.Run(nodes, cancellationToken);

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