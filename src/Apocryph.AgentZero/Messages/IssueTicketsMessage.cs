using System;
using System.Numerics;

namespace Apocryph.AgentZero.Messages
{
    public class IssueTicketsMessage
    {
        public Guid For { get; set; }
        public Guid Target { get; set; }
        public BigInteger Amount { get; set; }

        public IssueTicketsMessage(Guid @for, Guid target, BigInteger amount)
        {
            For = @for;
            Target = target;
            Amount = amount;
        }
    }
}