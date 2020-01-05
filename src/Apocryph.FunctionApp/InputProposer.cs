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
    public static class InputProposer
    {
        public class State
        {
            public ISet<(string, object)> PendingInputs { get; set; } = new HashSet<(string, object)>();
        }

        [FunctionName(nameof(InputProposer))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("committerStream")] IAsyncEnumerable<IHashed<IAgentStep>> committerStream,
            [PerperStream("commandExecutorStream")] IAsyncEnumerable<(string, object)> commandExecutorStream,
            [PerperStream("outputStream")] IAsyncCollector<AgentInput> outputStream)
        {
            var state = await context.FetchStateAsync<State>() ?? new State();

            await Task.WhenAll(
                commandExecutorStream.ForEachAsync(async agentInput =>
                {
                    state.PendingInputs.Add(agentInput);
                    await context.UpdateStateAsync(state);
                }, CancellationToken.None),

                committerStream.ForEachAsync(async step =>
                {
                    switch (step.Value)
                    {
                        case AgentInput input:
                            state.PendingInputs.Remove((input.Sender, input.Message));
                            await context.UpdateStateAsync(state);
                            break;
                        case AgentOutput output:
                            var (sender, message) = state.PendingInputs.First();
                            await outputStream.AddAsync(new AgentInput
                            {
                                Previous = step.Hash,
                                State = output.State,
                                Sender = sender,
                                Message = message,
                            });
                            break;
                    }
                }, CancellationToken.None));
        }
    }
}