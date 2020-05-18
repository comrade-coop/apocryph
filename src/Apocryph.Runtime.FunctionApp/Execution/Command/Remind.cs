using System;

namespace Apocryph.Runtime.FunctionApp.Execution.Command
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