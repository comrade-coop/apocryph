using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Agent;
using Apocryph.FunctionApp.Command;
using Apocryph.FunctionApp.Ipfs;
using Apocryph.FunctionApp.Model;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.FunctionApp
{
    public static class StepHashCollector
    {
        [FunctionName(nameof(StepHashCollector))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("inputStream")] IAsyncEnumerable<ISigned<object>> inputStream,
            [PerperStream("outputStream")] IAsyncCollector<Hash> outputStream)
        {
            await inputStream.ForEachAsync(async input =>
            {
                switch (input.Value)
                {
                    case Commit commit:
                        await outputStream.AddAsync(commit.For);
                        break;
                    case Vote vote:
                        await outputStream.AddAsync(vote.For);
                        break;
                    case Proposal proposal:
                        await outputStream.AddAsync(proposal.For);
                        break;
                }
            }, CancellationToken.None);
        }
    }
}