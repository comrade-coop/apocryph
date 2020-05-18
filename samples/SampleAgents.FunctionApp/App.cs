using System;
using System.Collections.Generic;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Testbed;
using Apocryph.Agent;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace SampleAgents.FunctionApp
{
    [Obsolete]
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
            await _testbed.Setup(context, "SampleAgents.FunctionApp.Agents.AgentOneWrapper.Run", "SampleAgents.FunctionApp.App.Runtime", "SampleAgents.FunctionApp.App.Monitor", cancellationToken);
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