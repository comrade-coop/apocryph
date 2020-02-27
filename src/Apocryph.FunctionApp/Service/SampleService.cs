using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Model;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.FunctionApp.Service
{
    public static class SampleService
    {
        public class State
        {
            public Dictionary<string, object> Data { get; set; } = new Dictionary<string, object>();
        }

        [FunctionName("service_sample")]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("commandStream")] IAsyncEnumerable<object> commandStream,
            [PerperStream("outputStream")] IAsyncCollector<(string, object)> outputStream)
        {
            var state = await context.FetchStateAsync<State>() ?? new State();

            await commandStream.ForEachAsync(async command =>
            {
                switch (command)
                {
                    case Tuple<string, object> storeCommand:
                        state.Data[storeCommand.Item1] = storeCommand.Item2;
                        await context.UpdateStateAsync(state);
                        break;
                    case Tuple<string> readCommand:
                        await outputStream.AddAsync((
                            "sample",
                            (readCommand.Item1, state.Data[readCommand.Item1])));
                        break;
                }
            }, CancellationToken.None);
        }
    }
}