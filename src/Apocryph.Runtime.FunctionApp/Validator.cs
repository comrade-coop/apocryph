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
        private HashSet<byte[]>? _pendingSetChainBlockMessages = new HashSet<byte[]>();
        private HashSet<object>? _pendingCommands;
        private IAsyncCollector<Message<Block>>? _output;
        private Guid _chainId;
        private Node _node;

        [FunctionName(nameof(Validator))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("assigner")] IAsyncEnumerable<(Node, Node[])> assigner,
            [Perper("chainId")] Guid chainId,
            [Perper("lastBlock")] Block lastBlock,
            [Perper("pendingCommands")] HashSet<object> pendingCommands,
            [PerperStream("gossips")] IAsyncEnumerable<Block> gossips,
            [PerperStream("queries")] IAsyncEnumerable<Query<Block>> queries,
            [PerperStream("output")] IAsyncCollector<Message<Block>> output,
            CancellationToken cancellationToken)
        {
            _chainId = chainId;
            _lastBlock = lastBlock;
            _pendingCommands = pendingCommands;
            _output = output;

            await Task.WhenAll(
                HandleAssigner(assigner),
                HandleIBC(gossips, cancellationToken),
                HandleQueries(context, queries, cancellationToken));
        }

        private async Task HandleAssigner(IAsyncEnumerable<(Node, Node[])> assigner)
        {
            await foreach (var (node, nodes) in assigner)
            {
                _node = node;
            }
        }

        private async Task<bool> Validate(PerperStreamContext context, Node node, Block block)
        {
            var _sawClaimRewardMessage = false;
            foreach (var inputCommand in block.InputCommands)
            {
                if (_chainId == Guid.Empty && inputCommand is Invoke invokation)
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

            var executor = new Executor(node?.ToString()!,
                async input => await context.CallWorkerAsync<(byte[]?, (string, object[])[], IDictionary<Guid, string[]>, IDictionary<Guid, string>)>("AgentWorker", new { input }, default));
            var (newStates, newCommands, newCapabilities) = await executor.Execute(
                _lastBlock!.States, block.InputCommands, _lastBlock!.Capabilities);

            return block.Equals(new Block(_chainId, block.ProposerAccount, newStates, block.InputCommands, newCommands, newCapabilities));
            // Validate historical blocks as per protocol
        }

        private async Task HandleIBC(IAsyncEnumerable<Block> ibc, CancellationToken cancellationToken)
        {
            var executor = new Executor(_node.Chain, default!);

            await foreach (var block in ibc.WithCancellation(cancellationToken))
            {
                foreach (var command in block.Commands)
                {
                    if (executor.FilterCommand(command, _lastBlock!.Capabilities))
                    {
                        _pendingCommands!.Add(command);
                    }
                }

                if (_chainId == Guid.Empty)
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

        private async Task HandleQueries(PerperStreamContext context, IAsyncEnumerable<Query<Block>> queries, CancellationToken cancellationToken)
        {
            await foreach (var query in queries.WithCancellation(cancellationToken))
            {
                if (query.Receiver != _node) continue;

                var block = query.Value;
                var valid = await Validate(context, _node, block);
                await _output!.AddAsync(new Message<Block>(block, valid ? MessageType.Valid : MessageType.Invalid), cancellationToken);
            }
        }
    }
}