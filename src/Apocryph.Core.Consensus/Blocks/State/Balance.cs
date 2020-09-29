using System;
using System.Numerics;

namespace Apocryph.Core.Consensus.Blocks.State
{
    public class Balance
    {
        public BigInteger Amount { get; private set; }

        public void Transfer(Balance? to, BigInteger amount)
        {
            if (amount < 0)
            {
                throw new Exception("Invalid amount");
            }
            if (Amount < amount)
            {
                throw new Exception("Insufficient funds");
            }

            Amount -= amount;
            if (to != null)
            {
                to.Amount += amount;
            }
        }
    }
}