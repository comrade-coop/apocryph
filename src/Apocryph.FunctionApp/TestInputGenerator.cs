using System;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Model;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.FunctionApp
{
    public static class TestInputGenerator
    {
        [FunctionName(nameof(TestInputGenerator))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("outputStream")] IAsyncCollector<AgentInput> outputStream)
        {
            await Task.Delay(TimeSpan.FromSeconds(1));

            await outputStream.AddAsync(new AgentInput
            {
                Previous = new Hash {Bytes = new byte[] {0}},
                State = new { },
                Sender = "",
                Message = new {init = true},
            });
        }
    }
}