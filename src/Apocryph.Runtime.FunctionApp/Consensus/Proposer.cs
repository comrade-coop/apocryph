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
        private bool _accepted;

        private Node? _node;
        private IAsyncCollector<Query<Block>>? _output;

        public Proposer()
        {
            _channel = Channel.CreateUnbounded<Query<Block>>();
            _accepted = false;
        }

        [FunctionName(nameof(Proposer))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("node")] Node node,
            [Perper("nodes")] Node[] nodes,
            [PerperStream("queries")] IAsyncEnumerable<Query<Block>> queries,
            [PerperStream("output")] IAsyncCollector<Query<Block>> output,
            CancellationToken cancellationToken)
        {
            _node = node;
            _output = output;

            var opinion = default(Block?);
            if (_node.IsProposer)
            {
                opinion = await Propose();
            }

            var snowball = new Snowball<Block>(_node, 100, 0.6, 3,
                SnowballSend, SnowballRespond, opinion);
            await Task.WhenAll(
                RunSnowball(snowball, nodes, cancellationToken),
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
                if (_accepted) break;

                if (query.Receiver == _node)
                {
                    await _channel.Writer.WriteAsync(query, cancellationToken);
                }
            }
        }

        private async Task RunSnowball(Snowball<Block> snowball, Node[] nodes, CancellationToken cancellationToken)
        {
            var acceptedProposal = await snowball.Run(nodes, cancellationToken);
            _accepted = true;

            await _output!.AddAsync(new Query<Block>(acceptedProposal, _node!, _node!, QueryVerb.Accept), cancellationToken);
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