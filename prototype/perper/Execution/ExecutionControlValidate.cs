using System.Collections.Generic;

namespace Apocryph.Execution
{
    public class ExecutionControlValidate : IExecutionControl
    {
        public int FromSnapshot;
        public string SnapshotIdentifier;
        public List<ExecutionMessage> Messages;
    }
}
