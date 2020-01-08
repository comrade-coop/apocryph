using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Model;
using Microsoft.Azure.WebJobs;
using Microsoft.Extensions.Logging;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.FunctionApp
{
    public static class ValidatorRuntimeOutput
    {
        public class State
        {
            public Dictionary<Hash, Dictionary<ValidatorKey, ValidatorSignature>> CommitSignatures { get; } = new Dictionary<Hash, Dictionary<ValidatorKey, ValidatorSignature>>();
        }

        [FunctionName(nameof(ValidatorRuntimeOutput))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("validatorFilterStream")] IAsyncEnumerable<ISigned<AgentOutput>> validatorFilterStream,
            [PerperStream("runtimeStream")] IAsyncEnumerable<AgentOutput> runtimeStream,
            [PerperStream("outputStream")] IAsyncCollector<AgentOutput> outputStream,
            ILogger logger)
        {
            var state = await context.FetchStateAsync<State>() ?? new State();

            await Task.WhenAll(
                validatorFilterStream.ForEachAsync(async expectedOutput =>
                {
                    try
                    {
                        state.CommitSignatures[expectedOutput.Value.Previous] = expectedOutput.Value.CommitSignatures;
                        await context.UpdateStateAsync(state);
                    }
                    catch (Exception e)
                    {
                        logger.LogError(e.ToString());
                    }
                }, CancellationToken.None),

                runtimeStream.ForEachAsync(async output =>
                {
                    try
                    {
                        output.CommitSignatures = state.CommitSignatures[output.Previous];

                        await outputStream.AddAsync(output);
                    }
                    catch (Exception e)
                    {
                        logger.LogError(e.ToString());
                    }
                }, CancellationToken.None));
        }
    }
}