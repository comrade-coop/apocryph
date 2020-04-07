using System.Collections.Generic;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Agents.Testbed;
using Apocryph.Agents.Testbed.Api;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Agent.FunctionApp
{
    public class App
    {
        private readonly Testbed _testbed;

        public App(Testbed testbed)
        {
            _testbed = testbed;
        }

        [FunctionName("Setup")]
        public async Task Setup(
            [PerperStreamTrigger(RunOnStartup = true)] PerperStreamContext context,
            CancellationToken cancellationToken)
        {
            await _testbed.Setup(context, "AgentOne", "Runtime", "Monitor", cancellationToken);
        }

        [FunctionName("Runtime")]
        public async Task Runtime(
            [PerperStreamTrigger] PerperStreamContext context,
            [Perper("agentDelegate")] string agentDelegate,
            [PerperStream("commands")] IAsyncEnumerable<AgentCommands> commands,
            CancellationToken cancellationToken)
        {
            await _testbed.Runtime(context, agentDelegate, commands, cancellationToken);
        }

        [FunctionName("Monitor")]
        public async Task Monitor(
            [PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("commands")] IAsyncEnumerable<AgentCommands> commands,
            CancellationToken cancellationToken)
        {
            await _testbed.Monitor(commands, cancellationToken);
        }
    }
}