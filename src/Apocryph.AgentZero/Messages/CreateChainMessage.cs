using System;
using System.Numerics;

namespace Apocryph.AgentZero.Messages
{
    public class CreateChainMessage
    {
        public Guid ChainId { get; set; }
        public BigInteger InitialBalance { get; set; }
        public byte[] InitialBlockId { get; set; }

        public CreateChainMessage(Guid chainId, BigInteger initialBalance, byte[] initialBlockId)
        {
            ChainId = chainId;
            InitialBalance = initialBalance;
            InitialBlockId = initialBlockId;
        }
    }
}