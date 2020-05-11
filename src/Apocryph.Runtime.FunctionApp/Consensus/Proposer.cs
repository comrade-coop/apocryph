using System;
using System.Collections.Generic;
using System.Linq;
using System.Runtime.CompilerServices;
using System.Threading;
using System.Threading.Channels;
using System.Threading.Tasks;
using Apocryph.Runtime.FunctionApp.Consensus.Core;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Runtime.FunctionApp.Consensus
{
    public class Proposer
    {
        private readonly Channel<Query<Block>> _channel;
        private bool _committed;

        private Node? _node;
        private IAsyncCollector<object>? _output;
        private Snowball<Block>? _snowball;

        public Proposer()
        {
            _channel = Channel.CreateUnbounded<Query<Block>>();
            _committed = false;
        }

        [FunctionName(nameof(Proposer))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("node")] Node node,
            [Perper("nodes")] Node[] nodes,
            [PerperStream("queries")] IAsyncEnumerable<Query<Block>> queries,
            [PerperStream("output")] IAsyncCollector<object> output,
            CancellationToken cancellationToken)
        {
            _node = node;
            _output = output;

            var opinion = default(Block?);
            if (_node.IsProposer)
            {
                opinion = await Propose();
            }

            _snowball = new Snowball<Block>(_node, 100, 0.6, 3,
                SnowballSend, SnowballRespond, opinion);
            await Task.WhenAll(
                RunSnowball(nodes, cancellationToken),
                HandleQueries(queries, cancellationToken));
        }

        private Task<Block> Propose()
        {
            throw new NotImplementedException();
        }

        private async Task HandleQueries(IAsyncEnumerable<Query<Block>> queries, CancellationToken cancellationToken)
        {
            await foreach (var query in queries.WithCancellation(cancellationToken))
            {
                if (_committed) break;

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

        private async Task RunSnowball(Node[] nodes, CancellationToken cancellationToken)
        {
            var committedProposal = await _snowball!.Run(nodes, cancellationToken);
            _committed = true;

            await _output!.AddAsync(new Message<Block>(committedProposal, MessageType.Committed), cancellationToken);
        }

        private async IAsyncEnumerable<Query<Block>> SnowballSend(Query<Block>[] queries, [EnumeratorCancellation] CancellationToken cancellationToken)
        {
            await Task.WhenAll(queries.Select(q => _output!.AddAsync(q, cancellationToken)));
            await foreach (var answer in _channel.Reader.ReadAllAsync(cancellationToken).Take(queries.Length).WithCancellation(cancellationToken))
            {
                yield return answer;
            }
        }

        private static Query<Block> SnowballRespond(Query<Block> query, Block? value, Block? opinion)
        {
            var result = opinion ??
                         (value is null || query.Value.ProposerStake > value.ProposerStake ? query.Value : value);
            return new Query<Block>(result, query.Receiver, query.Sender, QueryVerb.Response);
        }
    }
}