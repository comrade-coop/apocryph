using System;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Agent.Api;
using Apocryph.Agent.Worker;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace SampleAgents.FunctionApp.Agents
{
    public class AgentTwo : IAgent<object>
    {
        public void Setup(IContext<object> context)
        {
            context.RegisterInstance<IPingPongMessage, PingPongMessage>();
        }

        public Task Run(IContext<object> context, object message, Guid? reference)
        {
            switch (message)
            {
                case PingPongMessage pingPongMessage when pingPongMessage.AgentTwo == null:
                    context.Invoke(pingPongMessage.AgentTwo!.Value, context.CreateInstance<IPingPongMessage>(i =>
                    {
                        i.AgentOne = pingPongMessage.AgentOne;
                        i.AgentTwo = context.CreateReference(new[] {typeof(PingPongMessage)});
                    }));
                    break;
                case PingPongMessage pingPongMessage:
                    context.Invoke(pingPongMessage.AgentTwo!.Value, context.CreateInstance<IPingPongMessage>(i =>
                    {
                        i.AgentOne = pingPongMessage.AgentOne;
                        i.AgentTwo = pingPongMessage.AgentTwo;
                        i.Content = "Ping";
                    }));
                    break;
            }

            return Task.FromResult(context);
        }
    }

    public class AgentTwoWrapper
    {
        [FunctionName(nameof(AgentTwoWrapper))]
        public async Task<WorkerOutput> Run([PerperWorkerTrigger] PerperWorkerContext context,
            [Perper("input")] WorkerInput input, CancellationToken cancellationToken)
        {
            return await new Worker<object>(new AgentTwo()).Run(input);
        }
    }
}