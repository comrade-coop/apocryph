using System.Collections.Generic;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Agent;
using Apocryph.FunctionApp.Command;
using Apocryph.FunctionApp.Model;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Bindings;
using Perper.WebJobs.Extensions.Model;
using Perper.WebJobs.Extensions.Triggers;

namespace Apocryph.FunctionApp
{
    public static class Runtime
    {
        private class State
        {
            public object AgentState { get; set; }
        }

        [FunctionName("Runtime")]
        public static async Task Run([PerperStreamTrigger] IPerperStreamContext context,
            [PerperStream("validatorStream")] IPerperStream<IAgentStep> validatorStream,
            [PerperStream("committerStream")] IPerperStream<(IAgentStep, bool)> committerStream,
            [PerperStream] IAsyncCollector<(IAgentStep, IAgentStep)> outputStream)
        {
            await Task.WhenAll(
                validatorStream.Listen(async step =>
                {
                    switch (step)
                    {
                        case AgentOutput output:
                            //
                            break;
                        case AgentInput input:
                            //Call agent
                            await outputStream.AddAsync((
                                input,
                                new AgentOutput {Type = "Valid"}));
                            break;
                    }
                }, CancellationToken.None),

                committerStream.Listen(async commit =>
                {
                    var (step, isProposer) = commit;
                    switch (step)
                    {
                        case AgentOutput output:
                            var state = await context.GetState<State>("state");
                            state.AgentState = output.State;
                            await context.SetState("state", state);
                            if (isProposer)
                            {
                                foreach (var command in output.Commands)
                                {
                                    if (command is ReminderCommand reminder)
                                    {
                                        await Task.Delay(reminder.Time);
                                    }
                                }

                                await outputStream.AddAsync((
                                    output,
                                    new AgentInput {Type = "Proposal", State = state, Message = null, Sender = null}));
                            }
                            break;
                        case AgentInput input:
                            var result = new List<ICommand>();
                            var agentContext = new AgentContext<object>(result);
                            //Call agent

                            await outputStream.AddAsync((
                                input,
                                new AgentOutput {Type = "Proposal", Commands = result}));
                            break;
                    }
                }, CancellationToken.None));
        }
    }
}