using System;
using System.Collections.Generic;
using System.Numerics;
using System.Text.Json;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Core.Consensus;
using Apocryph.Core.Consensus.Blocks;
using Apocryph.Core.Consensus.Blocks.Command;
using Apocryph.Core.Consensus.Blocks.Messages;
using Apocryph.Core.Consensus.Communication;
using Apocryph.Core.Consensus.VirtualNodes;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Runtime.FunctionApp
{
    public class Validator
    {
        private Block? _lastBlock;
        private Dictionary<Block, Task<bool>> _validatedBlocks = new Dictionary<Block, Task<bool>>();
        private HashSet<byte[]>? _pendingSetChainBlockMessages = new HashSet<byte[]>();
        private HashSet<object>? _pendingCommands;
        private IAsyncCollector<Message<Block>>? _output;
        private Node? _node;

        [FunctionName(nameof(Validator))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("node")] Node node,
            [Perper("lastBlock")] Block lastBlock,
            [Perper("pendingCommands")] HashSet<object> pendingCommands,
            [PerperStream("consensus")] IAsyncEnumerable<Message<Block>> consensus,
            [PerperStream("filter")] IAsyncEnumerable<Block> filter,
            [PerperStream("queries")] IAsyncEnumerable<Query<Block>> queries,
            [PerperStream("output")] IAsyncCollector<Message<Block>> output,
            CancellationToken cancellationToken)
        {
            _lastBlock = lastBlock;
            _pendingCommands = pendingCommands;
            _output = output;
            _node = node;

            await Task.WhenAll(
                HandleFilter(filter, cancellationToken),
                HandleConsensus(context, consensus, cancellationToken),
                HandleQueries(context, queries, cancellationToken));
        }

        private async Task<bool> Validate(PerperStreamContext context, Node node, Block block)
        {
            var _sawClaimRewardMessage = false;
            foreach (var inputCommand in block.InputCommands)
            {
                if (_node!.ChainId.Length == 0 && inputCommand is Invoke invokation)
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

            var executor = new Executor(_node!.ChainId,
                async input => await context.CallWorkerAsync<(byte[]?, (string, object[])[], IDictionary<Guid, string[]>, IDictionary<Guid, string>)>("AgentWorker", new { input }, default));
            var (newStates, newCommands, newCapabilities) = await executor.Execute(
                _lastBlock!.States, block.InputCommands, _lastBlock!.Capabilities);

            return block.Equals(new Block(_node!.ChainId, block.ProposerAccount, newStates, block.InputCommands, newCommands, newCapabilities));
            // Validate historical blocks as per protocol
        }

        private async Task HandleFilter(IAsyncEnumerable<Block> filter, CancellationToken cancellationToken)
        {
            var executor = new Executor(_node!.ChainId, default!);

            await foreach (var block in filter.WithCancellation(cancellationToken))
            {
                foreach (var command in block.Commands)
                {
                    if (executor.FilterCommand(command, _lastBlock!.Capabilities))
                    {
                        _pendingCommands!.Add(command);
                    }
                }

                if (_node!.ChainId.Length == 0)
                {
                    _pendingSetChainBlockMessages!.Add(JsonSerializer.SerializeToUtf8Bytes(new SetChainBlockMessage
                    {
                        ChainId = block.ChainId,
                        BlockId = new byte[] { },
                        ProcessedCommands = new Dictionary<Guid, BigInteger>()
                        {
                            [block.ProposerAccount] = block.InputCommands.Length,
                        },
                        UsedTickets = new Dictionary<Guid, BigInteger>() { }, // TODO: Keep track of tickets
                        UnlockedTickets = new Dictionary<Guid, BigInteger>() { },
                    }));
                }
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
                if (query.Receiver != _node) continue;

                var block = query.Value;
                if (!_validatedBlocks.ContainsKey(block))
                {
                    _validatedBlocks[block] = Validate(context, _node, block);
                }
            }
        }
    }
}