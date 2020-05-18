using System;

namespace Apocryph.Agent.Command
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