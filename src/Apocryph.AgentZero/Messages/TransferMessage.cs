using System;
using System.Numerics;

namespace Apocryph.AgentZero.Messages
{
    public class TransferMessage
    {
        public Guid To { get; set; }
        public BigInteger Amount { get; set; }

        public TransferMessage(Guid to, BigInteger amount)
        {
            To = to;
            Amount = amount;
        }
    }
}