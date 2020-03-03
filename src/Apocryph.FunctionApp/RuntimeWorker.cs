using System;
using Apocryph.FunctionApp.Agent;
using Apocryph.FunctionApp.Service;
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
            if (message is InitMessage)
            {
                context.SampleStore("0", "1");
                context.AddReminder(TimeSpan.FromSeconds(5), "0");
            }
            else if (sender == "Reminder" && message is string key)
            {
                context.SampleRestore(key);
            }
            else if (sender == "Sample" && message is Tuple<string, object> data)
            {
                if (data.Item2 is string item)
                {
                    context.AddReminder(TimeSpan.FromSeconds(5), item);
                }
                else
                {
                    context.SampleStore(data.Item1, (int.Parse(data.Item1) + 1).ToString());
                    context.AddReminder(TimeSpan.FromSeconds(5), "0");
                }
            }
            return context;
        }
    }
}