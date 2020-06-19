using System;
using System.Collections.Generic;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Core.Consensus;
using Apocryph.Core.Consensus.Blocks;
using Apocryph.Core.Consensus.Communication;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Runtime.FunctionApp
{
    public class FilterStream
    {
        private Dictionary<byte[], Validator> _validators = new Dictionary<byte[], Validator>();
        private Dictionary<Block, Task<bool>> _validatedBlocks = new Dictionary<Block, Task<bool>>();
        private IAsyncCollector<Block>? _output;
        private Validator? _validator;

        [FunctionName(nameof(FilterStream))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("lastBlocks")] Dictionary<byte[], Block> lastBlocks,
            [PerperStream("pendingCommands")] Dictionary<byte[], HashSet<object>> pendingCommands,
            [PerperStream("ibc")] IAsyncEnumerable<Message<Block>> ibc,
            [PerperStream("gossips")] IAsyncEnumerable<Gossip<Block>> gossips,
            [PerperStream("output")] IAsyncCollector<Block> output,
            CancellationToken cancellationToken)
        {
            _output = output;

            foreach (var (chainId, lastBlock) in lastBlocks)
            {
                var executor = new Executor(chainId,
                    async input => await context.CallWorkerAsync<(byte[]?, (string, object[])[], IDictionary<Guid, string[]>, IDictionary<Guid, string>)>("AgentWorker", new { input }, default));
                _validators[chainId] = new Validator(executor, chainId, lastBlock, pendingCommands[chainId]);
            }

            await Task.WhenAll(
                HandleIBC(context, ibc, cancellationToken),
                HandleGossips(context, gossips, cancellationToken));
        }

        private Task<bool> Validate(PerperStreamContext context, Block block)
        {
            return _validators[block.ChainId].Validate(block);
        }

        private async Task HandleIBC(PerperStreamContext context, IAsyncEnumerable<Message<Block>> ibc, CancellationToken cancellationToken)
        {
            await foreach (var message in ibc.WithCancellation(cancellationToken))
            {
                if (message.Type != MessageType.Accepted) continue;

                var block = message.Value;
                if (!_validatedBlocks.ContainsKey(block))
                {
                    _validatedBlocks[block] = Validate(context, block);
                }

                var valid = await _validatedBlocks[block];
                if (valid)
                {
                    foreach (var (chainId, validator) in _validators)
                    {
                        validator.AddConfirmedBlock(block);
                    }

                    await _output!.AddAsync(block, cancellationToken);
                }
            }
        }

        private async Task HandleGossips(PerperStreamContext context, IAsyncEnumerable<Gossip<Block>> gossips, CancellationToken cancellationToken)
        {
            // Validate blocks from gossips before they are fully confirmed, saving a tiny bit of time
            await foreach (var gossip in gossips.WithCancellation(cancellationToken))
            {
                var block = gossip.Value;
                if (!_validatedBlocks.ContainsKey(block))
                {
                    _validatedBlocks[block] = Validate(context, block);
                }
            }
        }
    }
}