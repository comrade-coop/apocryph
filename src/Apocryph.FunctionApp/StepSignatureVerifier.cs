using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Agent;
using Apocryph.FunctionApp.Command;
using Apocryph.FunctionApp.Ipfs;
using Apocryph.FunctionApp.Model;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.FunctionApp
{
    public static class StepSignatureVerifier
    {
        [FunctionName(nameof(StepSignatureVerifier))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("stepsStream")] IAsyncEnumerable<IHashed<IAgentStep>> stepsStream,
            [Perper("validatorSet")] ValidatorSet validatorSet,
            [PerperStream("outputStream")] IAsyncCollector<Hash> outputStream)
        {
            await stepsStream.ForEachAsync(async step =>
            {
                bool validateCommit(ISigned<Commit> commit)
                {
                    if (commit.Value.For != step.Value.Previous)
                    {
                        return false;
                    }
                    var bytes = IpfsJsonSettings.ObjectToBytes(commit.Value);
                    return commit.Signer.ValidateSignature(bytes, commit.Signature);
                };

                if (step.Value.PreviousCommits.All(validateCommit))
                {
                    var committed = step.Value.PreviousCommits
                        .Select(commit => commit.Signer).Distinct()
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