using System;
using System.Collections.Generic;
using System.Linq;
using System.Numerics;
using System.Text;
using System.Text.Json;
using System.Threading.Tasks;
using Apocryph.Core.Consensus.VirtualNodes;
using Apocryph.Core.Consensus.Blocks;
using Apocryph.Core.Consensus.Blocks.Command;
using Apocryph.Core.Consensus.Blocks.Messages;

namespace Apocryph.Core.Consensus
{
    public class Proposer
    {
        private Guid _chainId;
        private Block _lastBlock;
        private Node? _proposer;
        private Guid _proposerAccount;
        private HashSet<Block> _confirmedBlocks;
        private HashSet<ICommand> _pendingCommands;
        private TaskCompletionSource<bool>? _pendingCommandsTaskCompletionSource;
        private Executor _executor;

        public Proposer(Executor executor, Guid chainId, Block lastBlock, HashSet<Block> confirmedBlocks, HashSet<ICommand> pendingCommands, Node? proposer, Guid proposerAccount)
        {
            _executor = executor;
            _chainId = chainId;
            _lastBlock = lastBlock;
            _confirmedBlocks = confirmedBlocks;
            _pendingCommands = pendingCommands;
            _proposer = proposer;
            _proposerAccount = proposerAccount;
        }


        public Block GetLastBlock()
        {
            return _lastBlock;
        }


        public async Task<Block> Propose()
        {
            if (_pendingCommands!.Count == 0)
            {
                _pendingCommandsTaskCompletionSource = new TaskCompletionSource<bool>();
                // TODO: Possible race condition if TrySetResult happens before assigning a new completion source
                await _pendingCommandsTaskCompletionSource.Task;
                _pendingCommandsTaskCompletionSource = null;
            }

            var inputCommands = _pendingCommands.ToArray();

            if (_chainId == Guid.Empty)
            {
                inputCommands = inputCommands.Concat(new ICommand[] {
                    new Invoke(_proposerAccount, (
                        "Apocryph.AgentZero.Messages.ClaimRewardMessage, Apocryph.AgentZero",
                        Encoding.UTF8.GetBytes("{}")))
                }).ToArray();
            }

            var (newState, newCommands, newCapabilities) = await _executor.Execute(
                _lastBlock!.States, inputCommands, _lastBlock.Capabilities);

            // Include historical blocks as per protocol
            var result = new Block(_chainId, _proposer, _proposerAccount, newState, inputCommands, newCommands, newCapabilities);

            _proposerAccount = Guid.NewGuid();

            return result;
        }


        public void AddConfirmedBlock(Block block)
        {
            if (!_confirmedBlocks.Add(block)) return;

            _pendingCommands!.UnionWith(block.Commands.Where(x => _executor.FilterCommand(x, _lastBlock!.Capabilities)));
            if (_pendingCommands!.Count > 0)
            {
                _pendingCommandsTaskCompletionSource?.TrySetResult(true);
            }

            if (_chainId == Guid.Empty)
            {
                _pendingCommands!.Add(new Invoke(_proposerAccount, (
                    typeof(SetChainBlockMessage).FullName!,
                    JsonSerializer.SerializeToUtf8Bytes(new SetChainBlockMessage
                    {
                        ChainId = block!.ChainId,
                        BlockId = new byte[] { },
                        ProcessedCommands = new Dictionary<Guid, BigInteger>()
                        {
                            [block.ProposerAccount] = block.InputCommands.Length,
                        },
                        UsedTickets = new Dictionary<Guid, BigInteger>() { }, // TODO: Keep track of tickets
                        UnlockedTickets = new Dictionary<Guid, BigInteger>() { },
                    }))));
            }

            if (_chainId == block.ChainId)
            {
                _lastBlock = block;
                _pendingCommands.ExceptWith(block.InputCommands);
            }
        }
    }
}