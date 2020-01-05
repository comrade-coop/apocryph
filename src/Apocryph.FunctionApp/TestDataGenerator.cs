using System;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Model;
using Microsoft.Azure.WebJobs;
using Microsoft.Extensions.Logging;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.FunctionApp
{
    public static class TestDataGenerator
    {
        [FunctionName(nameof(TestDataGenerator))]
        public static async Task RunAsync([PerperStreamTrigger] PerperStreamContext context,
            [Perper("delay")] TimeSpan delay,
            [Perper("data")] object data,
            [PerperStream("outputStream")] IAsyncCollector<object> outputStream,
            ILogger logger)
        {
            await Task.Delay(delay);
            logger.LogDebug("Sending a {type}", data.GetType());
            await outputStream.AddAsync(data);

            await context.BindOutput(CancellationToken.None);
        }
    }
}