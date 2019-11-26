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
    public static class Committer
    {
        private class State
        {
            public Dictionary<(AgentInput, AgentOutput), HashSet<string>> Votes { get; set; }
        }
        
        [FunctionName("Committer")]
        public static async Task Run([PerperStreamTrigger] IPerperStreamContext context,
            [PerperStream("validatorSet")] ValidatorSet validatorSet,
            [PerperStream("commitsStream")] IPerperStream<CommitMessage> commitsStream,
            [PerperStream] IAsyncCollector<(AgentOutput, bool)> outputStream)
        {
            var state = await context.GetState<State>("state");
            
            await commitsStream.Listen(
                async commit =>
                {
                    state.Votes[(commit.Input, commit.Output)].Add(commit.Signer);
                    await context.SetState("state", state);
                    
                    var voted = 0; //Count based on weights in validatorSet
                    if (3 * voted > 2 * validatorSet.Total)
                    {
                        const bool isProposer = true; //Check who is the next proposer

                        await outputStream.AddAsync((commit.Output, isProposer));
                    }
                },
                CancellationToken.None);
        }
    }
}