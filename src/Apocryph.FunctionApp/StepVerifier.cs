using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Agent;
using Apocryph.FunctionApp.Command;
using Apocryph.FunctionApp.Model;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.FunctionApp
{
    public static class StepVerifier
    {
        [FunctionName("StepVerifier")]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("stepsStream")] IAsyncEnumerable<Signed<IAgentStep>> stepsStream,
            [Perper("validatorSet")] ValidatorSet validatorSet,
            [PerperStream("outputStream")] IAsyncCollector<Hash> outputStream)
        {
            await stepsStream.ForEachAsync(async step =>
            {
                if (step.Value.CommitSignatures
                    .All(kv => kv.Key.ValidateSignature(step.Value.Previous, kv.Value)))
                {
                    var committed = step.Value.CommitSignatures.Keys
                        .Select(signer => validatorSet.Weights[signer]).Sum();
                    if (3 * committed > 2 * validatorSet.Total)
                    {
                        await outputStream.AddAsync(step.Value.Previous);
                    }
                }
            }, CancellationToken.None);
        }
    }
}