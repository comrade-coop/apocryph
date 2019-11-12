using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.Azure.WebJobs;
using Mock.Perper;
using Apocryph.Consensus;

public static class ConsensusPre
{
    public class State
    {
    }

    [FunctionName("ConsensusPre")]
    public static async Task RunAsync(PerperStreamContext<State> context,
        IAsyncEnumerable<ConsensusMessage> inputStream,
        [PerperOutput("consensus")] IAsyncCollector<object> postStream,
        [PerperOutput("execution")] IAsyncCollector<object> executionStream)
    {
        await inputStream.ForEachAsync(input =>
        {
            switch (input)
            {
                case VoteMessage v:
                    break;
            }
        });
    }
}

public static class ConsensusPost
{
    [FunctionName("ConsensusPre")]
    public static void Run(PerperStreamContext<object> context)
    {
    }
}
