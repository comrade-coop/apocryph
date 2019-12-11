using System.Collections.Generic;
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
        public static async Task Run([Perper(Stream = "Runtime")] IPerperStreamContext context,
            [Perper("validatorStream")] IAsyncEnumerable<IAgentStep> validatorStream,
            [Perper("committerStream")] IAsyncEnumerable<(IAgentStep, bool)> committerStream,
            [Perper("outputStream")] IAsyncCollector<(IAgentStep, bool)> outputStream)
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
                                    PreviousHash = input.Hash,
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
                                        PreviousHash = output.Hash,
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
                                        PreviousHash = input.Hash,
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