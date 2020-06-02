using System;

namespace Apocryph.Core.Consensus.Blocks.Command
{
    public class Remind
    {
        public DateTime DueDateTime { get; }
        public (string, byte[]) Message { get; }

        public Remind(DateTime dueDateTime, (string, byte[]) message)
        {
            DueDateTime = dueDateTime;
            Message = message;
        }
    }
}