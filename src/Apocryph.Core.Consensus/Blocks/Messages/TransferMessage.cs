using System;
using System.Numerics;

namespace Apocryph.Core.Consensus.Blocks.Messages
{
    public class TransferMessage
    {
        public Guid To { get; set; }
        public BigInteger Amount { get; set; }
    }
}