using System;

namespace Apocryph.Runtime.FunctionApp.Execution.Command
{
    public class Invoke
    {
        public Guid Reference { get; }
        public (string, byte[]) Message { get; }

        public Invoke(Guid reference, (string, byte[]) message)
        {
            Reference = reference;
            Message = message;
        }
    }
}