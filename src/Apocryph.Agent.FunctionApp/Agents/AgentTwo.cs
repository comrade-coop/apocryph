using System.Collections.Generic;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Agents.Testbed;
using Apocryph.Agents.Testbed.Api;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Agent.FunctionApp.Agents
{
    public class AgentTwo
    {
        public Task<AgentContext> Run(object state, AgentCapability self, object message)
        {
            var context = new AgentContext(state, self);
            if(message is PingPongMessage initMessage && initMessage.AgentTwo == null)
            {
                var cap = context.IssueCapability(new[] {typeof(PingPongMessage)});
                context.SendMessage(initMessage.AgentOne, new PingPongMessage
                {
                    AgentOne = initMessage.AgentOne,
                    AgentTwo = cap
                }, null);
            }
            else if(message is PingPongMessage pingPongMessage)
            {
                context.SendMessage(pingPongMessage.AgentOne, new PingPongMessage
                {
                    AgentOne = pingPongMessage.AgentOne,
                    AgentTwo = pingPongMessage.AgentTwo,
                    Content = "Pong"
                }, null);
            }
            return Task.FromResult(context);
        }
    }

    public class AgentTwoWrapper
    {
        private readonly Testbed _testbed;

        public AgentTwoWrapper(Testbed testbed)
        {
            _testbed = testbed;
        }

        [FunctionName("AgentTwo")]
        public async Task AgentTwo(
            [PerperStreamTrigger] PerperStreamContext context,
            [Perper("agentId")] string agentId,
            [Perper("initMessage")] object initMessage,
            [PerperStream("commands")] IAsyncEnumerable<AgentCommands> commands,
            [PerperStream("output")] IAsyncCollector<AgentCommands> output,
            CancellationToken cancellationToken)
        {
            await _testbed.Agent(new AgentTwo().Run, agentId, initMessage, commands, output, cancellationToken);
        }
    }
}