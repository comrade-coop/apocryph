using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Agent;
using Apocryph.FunctionApp.Command;
using Apocryph.FunctionApp.Model;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.FunctionApp
{
    public static class Runtime
    {
        [FunctionName("Runtime")]
        public static async Task Run([PerperTrigger("Runtime")] IPerperStreamContext context,
            [Perper("validatorStream")] IAsyncEnumerable<Hashed<IAgentStep>> validatorStream,
            [Perper("committerStream")] IAsyncEnumerable<(Hashed<IAgentStep>, bool)> committerStream,
            [Perper("outputStream")] IAsyncCollector<(IAgentStep, bool)> outputStream)
        {
            await Task.WhenAll(
                validatorStream.ForEachAsync(async step =>
                {
                    switch (step.Value)
                    {
                        case AgentOutput output:
                            // ..
                            break;
                        case AgentInput input:
                            var agentContext = await context.CallWorkerFunction<AgentContext<object>>(new
                            {
                                state = input.State,
                                sender = input.Sender,
                                message = input.Message
                            });
                            await outputStream.AddAsync((
                                new AgentOutput
                                {
                                    Previous = step.Hash,
                                    State = agentContext.State,
                                    Commands = agentContext.Commands
                                },
                                false));
                            break;
                    }
                }, CancellationToken.None),

                committerStream.ForEachAsync(async commit =>
                {
                    var (step, isProposer) = commit;
                    switch (step.Value)
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
                                        Previous = step.Hash,
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
                                var agentContext = await context.CallWorkerFunction<AgentContext<object>>(new
                                {
                                    state = input.State,
                                    sender = input.Sender,
                                    message = input.Message
                                });

                                await outputStream.AddAsync((
                                    new AgentOutput
                                    {
                                        Previous = step.Hash,
                                        State = agentContext.State,
                                        Commands = agentContext.Commands
                                    },
                                    true));
                            }
                            break;
                    }
                }, CancellationToken.None));
        }
    }
}