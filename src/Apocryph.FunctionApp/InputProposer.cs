using System;
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
            public (Hash, object)? PendingOutput { get; set; } = null;
        }

        [FunctionName(nameof(InputProposer))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("proposerFilterStream")] IAsyncEnumerable<IHashed<IAgentStep>> proposerFilterStream,
            [PerperStream("commandExecutorStream")] IAsyncEnumerable<(string, object)> commandExecutorStream,
            [PerperStream("outputStream")] IAsyncCollector<AgentInput> outputStream)
        {
            var state = await context.FetchStateAsync<State>() ?? new State();

            async Task processPending()
            {
                if (state.PendingInputs.Count != 0 && state.PendingOutput != null)
                {
                    var (sender, message) = state.PendingInputs.First();
                    var (previous, agentState) = state.PendingOutput.Value;

                    await outputStream.AddAsync(new AgentInput
                    {
                        Previous = previous,
                        State = agentState,
                        Sender = sender,
                        Message = message,
                    });

                    state.PendingOutput = null;
                }
            };

            await Task.WhenAll(
                commandExecutorStream.ForEachAsync(async agentInput =>
                {
                    state.PendingInputs.Add(agentInput);
                    await processPending();
                    await context.UpdateStateAsync(state);
                }, CancellationToken.None),

                proposerFilterStream.ForEachAsync(async step =>
                {
                    switch (step.Value)
                    {
                        case AgentInput input:
                            state.PendingInputs.Remove((input.Sender, input.Message));
                            state.PendingOutput = null;
                            break;
                        case AgentOutput output:
                            state.PendingOutput = (step.Hash, output.State);
                            break;
                    }
                    await processPending();
                    await context.UpdateStateAsync(state);
                }, CancellationToken.None));
        }
    }
}