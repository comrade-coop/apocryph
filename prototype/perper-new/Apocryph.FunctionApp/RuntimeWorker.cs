using System;
using Apocryph.FunctionApp.Agent;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.FunctionApp
{
    public static class RuntimeWorker
    {
        [FunctionName("RuntimeWorker")]
        public static void Run([Perper(Stream = "Runtime")] IPerperWorkerContext context,
            [Perper("agentContext", State = true)] IAgentContext<object> agentContext,
            [Perper("sender")] string sender,
            [Perper("message")] object message)
        {
            agentContext.AddReminder(TimeSpan.FromMinutes(5));
        }
    }
}