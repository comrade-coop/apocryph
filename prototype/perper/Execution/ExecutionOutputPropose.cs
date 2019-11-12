using System.Collections.Generic;

namespace Apocryph.Execution
{
    public class ExecutionOutputPropose : IExecutionOutput
    {
        public int Snapshot;
        public List<ExecutionMessage> Messages;
    }
}
