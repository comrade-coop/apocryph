using System;

namespace Apocryph.Execution
{
    public class ExecutionControlProposeStart : IExecutionControl
    {
        public int FromSnapshot;
        public int ChunkCountLimit = 1;
        public TimeSpan ChunkTimeLimit = TimeSpan.MaxValue;
    }
}
