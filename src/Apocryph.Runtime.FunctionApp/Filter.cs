using System;
using System.Collections.Generic;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Core.Consensus;
using Apocryph.Core.Consensus.Blocks;
using Apocryph.Core.Consensus.Blocks.Command;
using Apocryph.Core.Consensus.Blocks.Messages;
using Apocryph.Core.Consensus.Communication;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Runtime.FunctionApp
{
    public class Filter
    {
        private Block? _lastBlock;
        private Dictionary<Block, Task<bool>> _validatedBlocks = new Dictionary<Block, Task<bool>>();
        private HashSet<byte[]>? _pendingSetChainBlockMessages = new HashSet<byte[]>();
        private HashSet<object>? _pendingCommands;
        private IAsyncCollector<Block>? _output;

        [FunctionName(nameof(Filter))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("lastBlock")] Block lastBlock,
            [Perper("pendingCommands")] HashSet<object> pendingCommands,
            [PerperStream("ibc")] IAsyncEnumerable<Message<Block>> ibc,
            [PerperStream("gossips")] IAsyncEnumerable<Gossip<Block>> gossips,
            [PerperStream("output")] IAsyncCollector<Block> output,
            CancellationToken cancellationToken)
        {
            _lastBlock = lastBlock;
            _pendingCommands = pendingCommands;
            _output = output;

            await Task.WhenAll(
                HandleIBC(context, ibc, cancellationToken),
                HandleGossips(context, gossips, cancellationToken));
        }

        private async Task<bool> Validate(PerperStreamContext context, Block block)
        {
            var _sawClaimRewardMessage = false;
            foreach (var inputCommand in block.InputCommands)
            {
                if (block.ChainId.Length == 0 && inputCommand is Invoke invokation)
                {
                    if (invokation.Message.Item1 == typeof(ClaimRewardMessage).FullName)
                    {
                        if (_sawClaimRewardMessage)
                        {
                            return false;
                        }
                        _sawClaimRewardMessage = true;
                        continue;
                    }
                    else if (invokation.Message.Item1 == typeof(SetChainBlockMessage).FullName)
                    {
                        if (!_pendingSetChainBlockMessages!.Contains(invokation.Message.Item2))
                        {
                            return false;
                        }
                        continue;
                    }
                }
                if (!_pendingCommands!.Contains(inputCommand))
                {
                    return false;
                }
            }

            var executor = new Executor(block.ChainId,
                async input => await context.CallWorkerAsync<(byte[]?, (string, object[])[], IDictionary<Guid, string[]>, IDictionary<Guid, string>)>("AgentWorker", new { input }, default));
            var (newStates, newCommands, newCapabilities) = await executor.Execute(
                _lastBlock!.States, block.InputCommands, _lastBlock!.Capabilities);

            return block.Equals(new Block(block.ChainId, block.ProposerAccount, newStates, block.InputCommands, newCommands, newCapabilities));
            // Validate historical blocks as per protocol
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