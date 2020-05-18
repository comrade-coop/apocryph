using System;

namespace Apocryph.Agent.Worker
{
    public class WorkerInput
    {
        public byte[]? State { get; set; }
        public (string, byte[]) Message { get; }

        public Guid? Reference { get; set; }

        public WorkerInput((string, byte[]) message)
        {
            Message = message;
        }
    }
}