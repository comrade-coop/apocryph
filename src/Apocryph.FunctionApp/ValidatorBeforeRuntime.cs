using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Model;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.FunctionApp
{
    public static class ValidatorBeforeRuntime
    {
        public class State
        {
            public ISet<(string, object)> PendingInputs { get; set; } = new HashSet<(string, object)>();
        }

        [FunctionName("ValidatorBeforeRuntime")]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("validatorFilterStream")] IAsyncEnumerable<Signed<AgentOutput>> validatorFilterStream,
            [PerperStream("outputStream")] IAsyncCollector<Hash> outputStream)
        {
            var state = await context.FetchStateAsync<State>();

            await validatorFilterStream.ForEachAsync(async output =>
            {
                await outputStream.AddAsync(output.Value.Previous);
            }, CancellationToken.None);
        }
    }
}