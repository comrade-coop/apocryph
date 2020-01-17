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
        public static AgentContext<object> Run([PerperWorkerTrigger] object workerContext,
            [Perper("state")] object state,
            [Perper("sender")] string sender,
            [Perper("message")] object message)
        {
            var context = new AgentContext<object>(state);
            if (message is string stringMessage)
            {
                context.AddReminder(TimeSpan.FromSeconds(5), stringMessage + "i");
            }
            else
            {
                context.AddReminder(TimeSpan.FromSeconds(5), "h");
            }
            return context;
        }
    }
}