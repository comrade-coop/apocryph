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
    public static class Consensus
    {
        private class State
        {
            public Dictionary<(AgentInput, AgentOutput), HashSet<string>> Votes { get; set; }
        }
        
        [FunctionName("Consensus")]
        public static async Task Run([PerperStreamTrigger] IPerperStreamContext context,
            [PerperStream("validatorSet")] ValidatorSet validatorSet,
            [PerperStream("votesStream")] IPerperStream<Vote> votesStream,
            [PerperStream] IAsyncCollector<Commit> outputStream)
        {
            var state = await context.GetState<State>("state");
            
            await votesStream.Listen(
                async vote =>
                {
                    state.Votes[(vote.Input, vote.Output)].Add(vote.Signer);
                    await context.SetState("state", state);
                    
                    var voted = 0; //Count based on weights in validatorSet
                    if (3 * voted > 2 * validatorSet.Total)
                    {
                        await outputStream.AddAsync(new Commit {Input = vote.Input, Output = vote.Output});    
                    }
                },
                CancellationToken.None);
        }
    }
}