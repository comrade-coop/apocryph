using System;
using Apocryph.FunctionApp.Agent;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.FunctionApp
{
    public static class RuntimeWorker
    {
        [FunctionName(nameof(RuntimeWorker))]
        [return: Perper("$return")]
        public static AgentContext<object> Run([PerperWorkerTrigger("Runtime")] object state,
            [Perper("sender")] string sender,
            [Perper("message")] object message)
        {
            var context = new AgentContext<object>(state);
            context.AddReminder(TimeSpan.FromMinutes(5), new object {});
            return context;
        }
    }
}