using System.Collections.Generic;
using System.Threading;
using System.Threading.Tasks;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;
using SampleAgents.FunctionApp.Agents;

namespace SampleAgents.FunctionApp
{
    public static class Launcher
    {
        [FunctionName("Launcher")]
        public static async Task RunAsync([PerperStreamTrigger(RunOnStartup = true)]
            PerperStreamContext context,
            CancellationToken cancellationToken)
        {
            await context.StreamActionAsync("ChainList", new
            {
                chains = new Dictionary<byte[], string>
                {
                    {new byte[0], nameof(AgentOne)}
                }
            });

            await context.BindOutput(cancellationToken);
        }
    }
}