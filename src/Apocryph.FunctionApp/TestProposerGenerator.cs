using System.Threading.Tasks;
using Apocryph.FunctionApp.Model;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.FunctionApp
{
    public static class TestProposerGenerator
    {
        [FunctionName(nameof(TestProposerGenerator))]
        public static async Task RunAsync([PerperStreamTrigger] PerperStreamContext context,
            [Perper("self")] ValidatorKey self,
            [PerperStream("outputStream")] IAsyncCollector<ValidatorKey> outputStream)
        {
            await outputStream.AddAsync(self);
        }
    }
}