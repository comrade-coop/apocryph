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
        [return: Perper("$return")]
        public static AgentContext<object> Run([PerperWorker("Runtime")] object state,
            [Perper("sender")] string sender,
            [Perper("message")] object message)
        {
            var context = new AgentContext<object>(state);
            context.AddReminder(TimeSpan.FromMinutes(5));
            return context;
        }
    }
}