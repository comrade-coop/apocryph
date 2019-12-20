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
    public static class ValidatorRuntime
    {
        [FunctionName("ValidatorRuntime")]
        public static async Task Run([PerperStreamTrigger("ValidatorRuntime")] IPerperStreamContext context,
            [Perper("validatorStream")] IAsyncEnumerable<Hashed<IAgentStep>> validatorStream,
            [Perper("outputStream")] IAsyncCollector<IAgentStep> outputStream)
        {
            await validatorStream.ForEachAsync(async step =>
            {
                switch (step.Value)
                {
                    case AgentOutput output:
                        // ..
                        break;
                    case AgentInput input:
                        var agentContext = await context.CallWorkerAsync<AgentContext<object>>(new
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
                        break;
                }
            }, CancellationToken.None);
        }
    }
}