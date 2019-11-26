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
    public static class Voting
    {
        [FunctionName("Voting")]
        public static async Task Run([PerperStreamTrigger] IPerperStreamContext context,
            [PerperStream("contextStream")]
            IPerperStream<(AgentInput, AgentOutput)> contextStream,
            [PerperStream("proposalsStream")] 
            IPerperStream<(AgentInput, AgentOutput)> proposalsStream,
            [PerperStream] IAsyncCollector<object> outputStream)
        {
            var validOutputs = new Dictionary<AgentInput, AgentOutput>();
            await Task.WhenAll(
                contextStream.Listen(async validatedProposal =>
                {
                    var (validatedProposalInput, validatedProposalOutput) = validatedProposal;
                    if (validOutputs[validatedProposalInput] == validatedProposalOutput)
                    {
                        await outputStream.AddAsync(new VoteMessage
                            {Input = validatedProposalInput, Output = validatedProposalOutput});
                    }
                }, CancellationToken.None),
                proposalsStream.Listen(proposal =>
                {
                    var (proposedInput, proposedOutput) = proposal;
                    validOutputs[proposedInput] = proposedOutput;
                }, CancellationToken.None));
        }
    }
}