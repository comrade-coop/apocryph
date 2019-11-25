using System;
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
    public class Consensus
    {
        [FunctionName("Consensus")]
        public static async Task Run([PerperStreamTrigger] IPerperStreamContext context,
            [PerperStream("validStream")] IPerperStream<AgentOutput> validStream,
            [PerperStream("proposalsStream")] IPerperStream<(AgentInput, AgentOutput)> proposalsStream,
            [PerperStream] IAsyncCollector<object> outputStream)
        {
            var validOutputs = new Dictionary<AgentOutput, AgentInput>();
            await Task.WhenAll(
                validStream.Listen(async validOutput =>
                {
                    if (validOutputs.ContainsKey(validOutput))
                    {
                        await outputStream.AddAsync(new {vote = "agree"});
                    }
                }, CancellationToken.None),
                proposalsStream.Listen(proposal =>
                {
                    var (proposedInput, proposedOutput) = proposal;
                    validOutputs[proposedOutput] = proposedInput;
                }, CancellationToken.None));
        }
    }
}