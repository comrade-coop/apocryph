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
    public static class ProposerRuntime
    {
        [FunctionName("ProposerRuntime")]
        public static async Task Run([PerperTrigger("ProposerRuntime")] IPerperStreamContext context,
            [Perper("self")] ValidatorKey self,
            [Perper("committerStream")] IAsyncEnumerable<Hashed<IAgentStep>> committerStream,
            [Perper("currentProposerStream")] IAsyncEnumerable<ValidatorKey> currentProposerStream,
            [Perper("outputStream")] IAsyncCollector<IAgentStep> outputStream)
        {
            await committerStream.Zip(currentProposerStream).ForEachAsync(async pair =>
            {
                var (step, currentProposer) = pair;
                var isProposer = self.Equals(currentProposer);

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

                        if (isProposer) {
                            await outputStream.AddAsync(new AgentInput
                            {
                                Previous = step.Hash,
                                State = output.State,
                                Message = null, // ..
                                Sender = null, // ..
                            });
                        }
                        break;
                    case AgentInput input:
                        if (isProposer) {
                            var agentContext = await context.CallWorkerFunction<AgentContext<object>>(new
                            {
                                state = input.State,
                                sender = input.Sender,
                                message = input.Message
                            });

                            await outputStream.AddAsync(new AgentOutput
                            {
                                Previous = step.Hash,
                                State = agentContext.State,
                                Commands = agentContext.Commands
                            });
                        }
                        break;
                }
            }, CancellationToken.None);
        }
    }
}