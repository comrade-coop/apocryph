using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Model;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.FunctionApp
{
    public static class InputValidator
    {
        public class State
        {
            public ISet<(string, object)> PendingInputs { get; set; } = new HashSet<(string, object)>();
        }

        [FunctionName(nameof(InputValidator))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("validatorFilterStream")] IAsyncEnumerable<Signed<AgentInput>> validatorFilterStream,
            [PerperStream("committerStream")] IAsyncEnumerable<Hashed<AgentInput>> committerStream,
            [PerperStream("commandExecutorStream")] IAsyncEnumerable<(string, object)> commandExecutorStream,
            [PerperStream("outputStream")] IAsyncCollector<Vote> outputStream)
        {
            var state = await context.FetchStateAsync<State>() ?? new State();

            await Task.WhenAll(
                commandExecutorStream.ForEachAsync(async senderAndMessage =>
                {
                    state.PendingInputs.Add(senderAndMessage);
                    await context.UpdateStateAsync(state);
                }, CancellationToken.None),

                committerStream.ForEachAsync(async input =>
                {
                    state.PendingInputs.Remove((input.Value.Sender, input.Value.Message));
                    await context.UpdateStateAsync(state);
                }, CancellationToken.None),

                validatorFilterStream.ForEachAsync(async input =>
                {
                    var senderAndMessage = (input.Value.Sender, input.Value.Message);
                    if (state.PendingInputs.Contains(senderAndMessage))
                    {
                        await outputStream.AddAsync(new Vote { For = input.Hash });
                    }
                    // FIXME: Vote Nil?
                }, CancellationToken.None));
        }
    }
}