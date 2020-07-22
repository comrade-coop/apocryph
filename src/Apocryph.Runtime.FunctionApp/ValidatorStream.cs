using System;
using System.Collections.Generic;
using System.Threading;
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
    public class ValidatorStream
    {
        private Dictionary<Block, Task<bool>> _validatedBlocks = new Dictionary<Block, Task<bool>>();
        private IAsyncCollector<Message<Block>>? _output;
        private Node? _node;
        private Validator? _validator;

        [FunctionName(nameof(ValidatorStream))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("node")] Node node,
            [Perper("chainData")] Chain chainData,
            [Perper("consensus")] IAsyncEnumerable<Message<Block>> consensus,
            [Perper("filter")] IAsyncEnumerable<Block> filter,
            [Perper("queries")] IAsyncEnumerable<Query<Block>> queries,
            [Perper("output")] IAsyncCollector<Message<Block>> output,
            CancellationToken cancellationToken)
        {
            _output = output;
            _node = node;

            var executor = new Executor(_node!.ChainId,
                async (worker, input) => await context.CallWorkerAsync<(byte[]?, (string, object[])[], Dictionary<Guid, string[]>, Dictionary<Guid, string>)>(worker, new { input }, default));
            _validator = new Validator(executor, _node!.ChainId, chainData.GenesisBlock, new HashSet<object>());

            await TaskHelper.WhenAllOrFail(
                HandleFilter(filter, cancellationToken),
                HandleConsensus(context, consensus, cancellationToken),
                HandleQueries(context, queries, cancellationToken));
        }

        private Task<bool> Validate(PerperStreamContext context, Node node, Block block)
        {
            return _validator!.Validate(block);
            // Validate historical blocks as per protocol
        }

        private async Task HandleFilter(IAsyncEnumerable<Block> filter, CancellationToken cancellationToken)
        {
            await foreach (var block in filter.WithCancellation(cancellationToken))
            {
                _validator!.AddConfirmedBlock(block);
            }
        }


        private async Task HandleConsensus(PerperStreamContext context, IAsyncEnumerable<Message<Block>> consensus, CancellationToken cancellationToken)
        {
            await foreach (var message in consensus.WithCancellation(cancellationToken))
            {
                if (message.Type != MessageType.Proposed) continue;

                var block = message.Value;
                if (!_validatedBlocks.ContainsKey(block))
                {
                    _validatedBlocks[block] = Validate(context, _node!, block);
                }

                var valid = await _validatedBlocks[block];

                await _output!.AddAsync(new Message<Block>(block, valid ? MessageType.Valid : MessageType.Invalid), cancellationToken);
            }
        }

        private async Task HandleQueries(PerperStreamContext context, IAsyncEnumerable<Query<Block>> queries, CancellationToken cancellationToken)
        {
            // Validate blocks from queries before they are fully confirmed, saving a tiny bit of time
            await foreach (var query in queries.WithCancellation(cancellationToken))
            {
                if (!query.Receiver.Equals(_node)) continue;

                var block = query.Value;
                if (!_validatedBlocks.ContainsKey(block))
                {
                    _validatedBlocks[block] = Validate(context, _node, block);
                }
            }
        }
    }
}