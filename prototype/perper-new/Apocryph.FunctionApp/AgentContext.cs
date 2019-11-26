using System.Threading;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Model;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Bindings;
using Perper.WebJobs.Extensions.Model;
using Perper.WebJobs.Extensions.Triggers;

namespace Apocryph.FunctionApp
{
    public static class AgentContext
    {
        private class State
        {
            public object AgentState { get; set; }
        }

        [FunctionName("AgentContext")]
        public static async Task Run([PerperStreamTrigger] IPerperStreamContext context,
            [PerperStream("validatorStream")] IPerperStream<AgentInput> validatorStream,
            [PerperStream("committerStream")] IPerperStream<(AgentOutput, bool)> committerStream,
            [PerperStream] IAsyncCollector<(AgentInput, AgentOutput)> outputStream)
        {
            await Task.WhenAll(
                validatorStream.Listen(async input =>
                {
                    //Call agent
                    await outputStream.AddAsync((
                        input,
                        new AgentOutput {Type = "Valid"}));
                }, CancellationToken.None),

                committerStream.Listen(async commit =>
                {
                    var state = await context.GetState<State>("state");
                    state.AgentState = commit.Item1.State;
                    await context.SetState("state", state);

                    if (commit.Item2)
                    {
                        // Execute commands in the commitMessage
                        await Task.Delay(1000);

                        //Call agent
                        await outputStream.AddAsync((
                            new AgentInput {State = state.AgentState},
                            new AgentOutput {Type = "Proposal"}));
                    }
                }, CancellationToken.None));
        }
    }
}