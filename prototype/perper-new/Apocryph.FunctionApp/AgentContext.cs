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
            [PerperStream("commitStream")] IPerperStream<VoteMessage> commitStream,
            [PerperStream] IAsyncCollector<object> outputStream)
        {
            await Task.WhenAll(
                validatorStream.Listen(async validationMessage =>
                {
                    //Call Agent
                    await outputStream.AddAsync(new AgentOutput {Type = "Valid"});
                }, CancellationToken.None),

                commitStream.Listen(async commitMessage =>
                {
                    var state = await context.GetState<State>("state");
                    state.AgentState = commitMessage.Output.State;
                    await context.SetState("state", state);

                    // Execute commands in the commitMessage
                    await outputStream.AddAsync(new AgentOutput {Type = "Proposal"});
                }, CancellationToken.None));
        }
    }
}