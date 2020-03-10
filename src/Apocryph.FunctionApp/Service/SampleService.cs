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
        private class State
        {
            public Dictionary<string, object> Data { get; set; } = new Dictionary<string, object>();
        }

        [FunctionName("ServiceSample")]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("commandsStream")] IAsyncEnumerable<object> commandsStream,
            [PerperStream("outputStream")] IAsyncCollector<(string, object)> outputStream)
        {
            var state = await context.FetchStateAsync<State>() ?? new State();

            await commandsStream.ForEachAsync(async command =>
            {
                switch (command)
                {
                    case Tuple<string, object> storeCommand:
                        state.Data[storeCommand.Item1] = storeCommand.Item2;
                        await context.UpdateStateAsync(state);
                        break;
                    case Tuple<string> readCommand:
                        state.Data.TryGetValue(readCommand.Item1, out object value);
                        await outputStream.AddAsync((
                            "Sample",
                            Tuple.Create(readCommand.Item1, value)));
                        break;
                }
            }, CancellationToken.None);
        }
    }
}