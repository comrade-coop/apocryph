using System;
using System.Numerics;

namespace Apocryph.Core.Consensus.Blocks.Messages
{
    public class CreateChainMessage
    {
        public Guid ChainId { get; set; }
        public BigInteger InitialBalance { get; set; }
        public byte[] InitialBlockId { get; set; }
    }
}