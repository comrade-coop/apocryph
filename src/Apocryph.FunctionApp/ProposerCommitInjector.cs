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

    public static class ProposerCommitInjector
    {
        private class State
        {
            public Dictionary<Hash, List<ISigned<Commit>>> Commits { get; } = new Dictionary<Hash, List<ISigned<Commit>>>();
        }

        [FunctionName(nameof(ProposerCommitInjector))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("commitsStream")] IAsyncEnumerable<ISigned<Commit>> commitsStream,
            [PerperStream("stepsStream")] IAsyncEnumerable<IAgentStep> stepsStream,
            [PerperStream("outputStream")] IAsyncCollector<IAgentStep> outputStream,
            ILogger logger)
        {
            var state = await context.FetchStateAsync<State>() ?? new State();

            await Task.WhenAll(
                commitsStream.ForEachAsync(async commit =>
                {
                    try
                    {
                        if (!state.Commits.ContainsKey(commit.Value.For))
                        {
                            state.Commits[commit.Value.For] = new List<ISigned<Commit>>();
                        }
                        state.Commits[commit.Value.For].Add(commit);
                        await context.UpdateStateAsync(state);
                    }
                    catch (Exception e)
                    {
                        logger.LogError(e.ToString());
                    }
                }, CancellationToken.None),

                stepsStream.ForEachAsync(async step =>
                {
                    try
                    {
                        if (state.Commits.ContainsKey(step.Previous))
                        {
                            // FIXME: Should probably block until there are enough signatures (as we cannot be sure that the other stream would collect them in time)
                            step.PreviousCommits = state.Commits[step.Previous];
                        }

                        logger.LogDebug("Proposing a step after {0}!", step.Previous.Bytes[0]);

                        await outputStream.AddAsync(step);
                    }
                    catch (Exception e)
                    {
                        logger.LogError(e.ToString());
                    }
                }, CancellationToken.None));
        }
    }
}