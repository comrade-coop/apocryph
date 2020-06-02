using System;
using System.Numerics;

namespace Apocryph.Core.Consensus.Blocks.Messages
{
    public class IssueTicketsMessage
    {
        public Guid For { get; set; }
        public Guid Target { get; set; }
        public BigInteger Amount { get; set; }
    }
}