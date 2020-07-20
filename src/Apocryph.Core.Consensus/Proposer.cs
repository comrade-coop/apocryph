using System;
using System.Collections.Generic;
using System.Linq;
using System.Numerics;
using System.Text;
using System.Text.Json;
using System.Threading.Tasks;
using Apocryph.Core.Consensus.Blocks;
using Apocryph.Core.Consensus.Blocks.Command;
using Apocryph.Core.Consensus.Blocks.Messages;

namespace Apocryph.Core.Consensus
{
    public class Proposer
    {
        private Guid _chainId;
        private Block _lastBlock;
        private Guid _proposerAccount;
        private HashSet<object> _pendingCommands;
        private TaskCompletionSource<bool>? _pendingCommandsTaskCompletionSource;
        private Executor _executor;

        public Proposer(Executor executor, Guid chainId, Block lastBlock, HashSet<object> pendingCommands, Guid proposerAccount)
        {
            _executor = executor;
            _chainId = chainId;
            _lastBlock = lastBlock;
            _pendingCommands = pendingCommands;
            _proposerAccount = proposerAccount;
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
            _pendingCommands.Clear();

            if (_chainId == Guid.Empty)
            {
                inputCommands = inputCommands.Concat(new object[] {
                    new Invoke(_proposerAccount, (
                        "Apocryph.AgentZero.Messages.ClaimRewardMessage, Apocryph.AgentZero",
                        Encoding.UTF8.GetBytes("{}")))
                }).ToArray();
            }

            var (newState, newCommands, newCapabilities) = await _executor.Execute(
                _lastBlock!.States, inputCommands, _lastBlock.Capabilities);
            // Include historical blocks as per protocol
            return new Block(_chainId, _proposerAccount, newState, inputCommands, newCommands, newCapabilities);
        }


        public void AddConfirmedBlock(Block block)
        {

            _pendingCommandsTaskCompletionSource?.TrySetResult(true);
            _pendingCommands!.UnionWith(block.Commands.Where(x => _executor.FilterCommand(x, _lastBlock!.Capabilities)));

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