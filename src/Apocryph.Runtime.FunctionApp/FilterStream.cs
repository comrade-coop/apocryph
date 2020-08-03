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
        private Dictionary<Guid, Validator> _validators = new Dictionary<Guid, Validator>();
        private Dictionary<Block, Task<bool>> _validatedBlocks = new Dictionary<Block, Task<bool>>();
        private IAsyncCollector<Block>? _output;

        [FunctionName(nameof(FilterStream))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("chains")] Dictionary<Guid, Chain> chains,
            [Perper("ibc")] IAsyncEnumerable<Message<Block>> ibc,
            [Perper("gossips")] IAsyncEnumerable<Gossip<Block>> gossips,
            [Perper("output")] IAsyncCollector<Block> output,
            CancellationToken cancellationToken)
        {
            _output = output;

            foreach (var (chainId, chain) in chains)
            {
                var executor = new Executor(chainId,
                    async (worker, input) => await context.CallWorkerAsync<(byte[]?, (string, object[])[], Dictionary<Guid, string[]>, Dictionary<Guid, string>)>(worker, new { input }, default));
                _validators[chainId] = new Validator(executor, chainId, chain.GenesisBlock, new HashSet<Block>(), new HashSet<object>());
            }

            await TaskHelper.WhenAllOrFail(
                Task.Run(async () =>
                {

                    // Second loop, as we want to distribute all genesis blocks to all chains
                    foreach (var (chainId, chain) in chains)
                    {
                        await Task.Delay(2000);
                        foreach (var (_chainId, validator) in _validators)
                        {
                            validator.AddConfirmedBlock(chain.GenesisBlock);
                        }

                        await _output!.AddAsync(chain.GenesisBlock, cancellationToken);
                    }
                }),
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