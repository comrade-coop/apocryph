using System;
using System.Collections.Generic;
using System.Numerics;

namespace Apocryph.AgentZero.Messages
{
    public class SetChainBlockMessage
    {
        public Guid ChainId { get; set; }
        public byte[] BlockId { get; set; }

        public Dictionary<Guid, BigInteger> ProcessedCommands { get; set; } // Proposer => Amount
        public Dictionary<Guid, BigInteger> UsedTickets { get; set; } // Other chain => Tickets
        public Dictionary<Guid, BigInteger> UnlockedTickets { get; set; } // Other chain => Tickets

        public SetChainBlockMessage(
            Guid chainId,
            byte[] blockId,
            Dictionary<Guid, BigInteger> processedCommands,
            Dictionary<Guid, BigInteger> usedTickets,
            Dictionary<Guid, BigInteger> unlockedTickets)
        {
            ChainId = chainId;
            BlockId = blockId;
            ProcessedCommands = processedCommands;
            UsedTickets = usedTickets;
            UnlockedTickets = unlockedTickets;
        }
    }
}