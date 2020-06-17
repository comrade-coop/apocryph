using System;
using System.Collections.Generic;
using System.Linq;
using System.Numerics;
using System.Text.Json;
using System.Threading.Tasks;
using Apocryph.Core.Consensus.Blocks;
using Apocryph.Core.Consensus.Blocks.Command;
using Apocryph.Core.Consensus.Blocks.Messages;

namespace Apocryph.Core.Consensus
{
    public class Validator
    {
        private byte[] _chainId;
        private Block _lastBlock;
        private HashSet<object> _pendingCommands;
        private HashSet<byte[]>? _pendingSetChainBlockMessages = new HashSet<byte[]>();
        private Executor _executor;

        public Validator(Executor executor, byte[] chainId, Block lastBlock, HashSet<object> pendingCommands)
        {
            _executor = executor;
            _chainId = chainId;
            _lastBlock = lastBlock;
            _pendingCommands = pendingCommands;
        }

        public async Task<bool> Validate(Block block)
        {
            var _sawClaimRewardMessage = false;
            foreach (var inputCommand in block.InputCommands)
            {
                if (_chainId.Length == 0 && inputCommand is Invoke invokation)
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

            var inputCommands = _pendingCommands.ToArray();
            _pendingCommands.Clear();

            var (newState, newCommands, newCapabilities) = await _executor.Execute(
                _lastBlock!.States, block.InputCommands, _lastBlock!.Capabilities);

            // Validate historical blocks as per protocol
            return block.Equals(new Block(_chainId, block.ProposerAccount, newState, block.InputCommands, newCommands, newCapabilities));
        }


        public void AddConfirmedBlock(Block block)
        {

            _pendingCommands!.UnionWith(block.Commands.Where(x => _executor.FilterCommand(x, _lastBlock!.Capabilities)));

            if (_chainId.Length == 0)
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

            if (_chainId.SequenceEqual(block.ChainId))
            {
                _lastBlock = block;
                _pendingCommands.ExceptWith(block.InputCommands);
            }
        }
    }
}