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
    public static class StepOrderVerifier
    {
        private class State
        {
            public Hash CurrentStep { get; set; } = new Hash {Bytes = new byte[] {0}};
        }

        [FunctionName(nameof(StepOrderVerifier))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("stepsStream")] IAsyncEnumerable<IHashed<IAgentStep>> stepsStream,
            [PerperStream("outputStream")] IAsyncCollector<IHashed<IAgentStep>> outputStream)
        {
            var state = await context.FetchStateAsync<State>() ?? new State();

            await stepsStream.ForEachAsync(async step =>
            {
                if (step.Value.Previous == state.CurrentStep)
                {
                    state.CurrentStep = step.Hash;
                    await outputStream.AddAsync(step);
                    await context.UpdateStateAsync(state);
                }
            }, CancellationToken.None);
        }
    }
}