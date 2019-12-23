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
        public static async Task Run([PerperStreamTrigger("Runtime")] IPerperStreamContext context,
            [Perper("self")] ValidatorKey self,
            [Perper("inputStream")] IAsyncEnumerable<Hashed<AgentInput>> inputStream,
            [Perper("outputStream")] IAsyncCollector<AgentOutput> outputStream)
        {
            await inputStream.ForEachAsync(async input =>
            {
                var agentContext = await context.CallWorkerAsync<AgentContext<object>>(new
                {
                    state = input.Value.State,
                    sender = input.Value.Sender,
                    message = input.Value.Message
                });

                await outputStream.AddAsync(new AgentOutput
                {
                    Previous = input.Hash,
                    State = agentContext.State,
                    Commands = agentContext.Commands
                });
            }, CancellationToken.None);
        }
    }
}