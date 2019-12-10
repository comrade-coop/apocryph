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
        [FunctionName("Runtime")]
        public static async Task Run([PerperStreamTrigger] IPerperStreamContext context,
            [PerperStream("validatorStream")] IAsyncEnumerable<IAgentStep> validatorStream,
            [PerperStream("committerStream")] IAsyncEnumerable<(IAgentStep, bool)> committerStream,
            [PerperStream] IAsyncCollector<(IAgentStep, bool)> outputStream)
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
                            var result = new List<ICommand>();
                            var agentContext = new AgentContext<object>(result);
                            //Call agent
                            await outputStream.AddAsync((
                                new AgentOutput
                                {
                                    Previous = input,
                                    State = agentContext.State,
                                    Commands = result,
                                },
                                false));
                            break;
                    }
                }, CancellationToken.None),

                committerStream.Listen(async commit =>
                {
                    var (step, isProposer) = commit;
                    switch (step)
                    {
                        case AgentOutput output:
                            foreach (var command in output.Commands)
                            {
                                // execute services
                                if (command is ReminderCommand reminder)
                                {
                                    await Task.Delay(reminder.Time); // should not block other things executing
                                }
                            }
                            if (isProposer)
                            {
                                await outputStream.AddAsync((
                                    new AgentInput
                                    {
                                        Previous = output,
                                        State = output.State,
                                        Message = null, // ..
                                        Sender = null, // ..
                                    },
                                    true));
                            }
                            break;
                        case AgentInput input:
                            if (isProposer)
                            {
                                var result = new List<ICommand>();
                                var agentContext = new AgentContext<object>(result);
                                //Call agent
                                await outputStream.AddAsync((
                                    new AgentOutput
                                    {
                                        Previous = input,
                                        State = agentContext.State,
                                        Commands = result,
                                    },
                                    true));
                            }
                            break;
                    }
                }, CancellationToken.None));
        }
    }
}