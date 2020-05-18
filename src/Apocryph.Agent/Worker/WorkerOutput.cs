using System;
using System.Collections.Generic;

namespace Apocryph.Agent.Worker
{
    public class WorkerOutput
    {
        public byte[]? State { get; }
        public (string, object[])[] Actions { get; }
        public IDictionary<Guid, string[]> CreatedReferences { get; }
        public IDictionary<Guid, string> AttachedReferences { get; }

        public WorkerOutput(byte[]? state, (string, object[])[] actions,
            IDictionary<Guid, string[]> createdReferences,
            IDictionary<Guid, string> attachedReferences)
        {
            State = state;
            Actions = actions;
            CreatedReferences = createdReferences;
            AttachedReferences = attachedReferences;
        }
    }
}