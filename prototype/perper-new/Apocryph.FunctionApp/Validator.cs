using System.Collections.Generic;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Model;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Bindings;
using Perper.WebJobs.Extensions.Model;
using Perper.WebJobs.Extensions.Triggers;

namespace Apocryph.FunctionApp
{
    public static class Validator
    {
        private class State
        {
            public bool IsProposer { get; set; }
            public Dictionary<(AgentInput, AgentOutput), HashSet<string>> Votes { get; set; }
        }

        [FunctionName("Validator")]
        public static async Task Run([PerperStreamTrigger] IPerperStreamContext context,
            [PerperStream("validatorSet")] ValidatorSet validatorSet,
            [PerperStream("commitsStream")] IPerperStream<Commit> commitsStream,
            [PerperStream("proposalsStream")] IPerperStream<(AgentInput, AgentOutput)> proposalsStream,
            [PerperStream] IAsyncCollector<object> outputStream)
        {
            var state = await context.GetState<State>("state");

            await Task.WhenAll(
                commitsStream.Listen(async commit =>
                {
                    state.Votes[(commit.Input, commit.Output)].Add(commit.Signer);
                    
                    var voted = 0; //Count based on weights in validatorSet
                    if (3 * voted > 2 * validatorSet.Total)
                    {
                        state.IsProposer = true; //Check who is the next proposer
                    }
                    
                    await context.SetState("state", state);
                }, CancellationToken.None),

                proposalsStream.Listen(async proposal =>
                {
                    if (!state.IsProposer)
                    {
                        await outputStream.AddAsync(proposal.Item1);
                    }
                }, CancellationToken.None));
        }
    }
}