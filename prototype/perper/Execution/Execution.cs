using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.Azure.WebJobs;
using Mock.Perper;
using Apocryph.Execution;

namespace Apocryph.Execution
{
    public static class ExecutionPre
    {
        public class State
        {
        }

        [FunctionName("ExecutionPre")]
        public static async Task RunAsync(PerperStreamContext<State> context,
            IAsyncEnumerable<IExecutionControl> controlStream,
            [PerperOutput("messaging")] IAsyncCollector<(string, ExecutionMessage)> messagingStream)
        {
            await controlStream.ForEachAsync(async control =>
            {
                switch (control)
                {
                    case ExecutionControlProposeStart start:
                        break;
                    case ExecutionControlProposeEnd end:
                        break;
                    case ExecutionControlValidate validate:
                        break;
                    default:
                        // Assume we never get here
                        break;
                }
            });
        }
    }

    public static class ExecutionPost
    {
        public class State
        {
        }

        [FunctionName("ExecutionPost")]
        public static async Task RunAsync(PerperStreamContext<State> context,
            IAsyncEnumerable<(string,object)> executionOutputStream,
            IAsyncCollector<IExecutionOutput> outputStream)
        {
            await executionOutputStream.ForEachAsync(async output =>
            {
            });
        }
    }
}
