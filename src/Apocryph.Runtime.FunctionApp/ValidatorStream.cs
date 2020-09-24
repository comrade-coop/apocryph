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
    public class ValidatorStream
    {
        private Dictionary<Hash, Task<bool>> _validatedBlocks = new Dictionary<Hash, Task<bool>>();
        private IQueryable<HashRegistryEntry>? _hashRegistry;
        private IAsyncCollector<Message<Hash>>? _output;
        private Node? _node;
        private Validator? _validator;

        [FunctionName(nameof(ValidatorStream))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("node")] Node node,
            [Perper("chainData")] Chain chainData,
            [Perper("consensus")] IAsyncEnumerable<Message<Hash>> consensus,
            [Perper("filter")] IAsyncEnumerable<Hash> filter,
            [Perper("queries")] IAsyncEnumerable<Query<Hash>> queries,
            [Perper("hashRegistry")] IPerperStream hashRegistry,
            [Perper("output")] IAsyncCollector<Message<Hash>> output,
            CancellationToken cancellationToken)
        {
            _output = output;
            _node = node;
            _hashRegistry = context.Query<HashRegistryEntry>(hashRegistry);

            var executor = new Executor(_node!.ChainId,
                async (worker, input) => await context.CallWorkerAsync<(byte[]?, (string, object[])[], Dictionary<Guid, string[]>, Dictionary<Guid, string>)>(worker, new { input }, default));
            _validator = new Validator(executor, _node!.ChainId, chainData.GenesisBlock, new HashSet<Hash>(), new HashSet<ICommand>());

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

        private async Task HandleFilter(IAsyncEnumerable<Hash> filter, CancellationToken cancellationToken)
        {
            await foreach (var hash in filter.WithCancellation(cancellationToken))
            {
                var block = HashRegistryStream.GetObjectByHash<Block>(_hashRegistry!, hash);
                _validator!.AddConfirmedBlock(block!);
            }
        }


        private async Task HandleConsensus(PerperStreamContext context, IAsyncEnumerable<Message<Hash>> consensus, CancellationToken cancellationToken)
        {
            await foreach (var message in consensus.WithCancellation(cancellationToken))
            {
                if (message.Type != MessageType.Proposed) continue;

                var hash = message.Value;
                var block = HashRegistryStream.GetObjectByHash<Block>(_hashRegistry!, hash);
                if (!_validatedBlocks.ContainsKey(hash))
                {
                    _validatedBlocks[hash] = Validate(context, _node!, block!);
                }

                var valid = await _validatedBlocks[hash];

                if (valid)
                {
                    _validator!.AddConfirmedBlock(block!);
                }

                await _output!.AddAsync(new Message<Hash>(hash, valid ? MessageType.Valid : MessageType.Invalid), cancellationToken);
            }
        }

        private async Task HandleQueries(PerperStreamContext context, IAsyncEnumerable<Query<Hash>> queries, CancellationToken cancellationToken)
        {
            // Validate blocks from queries before they are fully confirmed, saving a tiny bit of time
            await foreach (var query in queries.WithCancellation(cancellationToken))
            {
                if (!query.Receiver.Equals(_node)) continue;

                var hash = query.Value;
                if (!_validatedBlocks.ContainsKey(hash))
                {
                    var block = HashRegistryStream.GetObjectByHash<Block>(_hashRegistry!, hash);
                    _validatedBlocks[hash] = Validate(context, _node, block!);
                }
            }
        }
    }
}