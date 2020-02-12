using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Agent;
using Apocryph.FunctionApp.Command;
using Apocryph.FunctionApp.Ipfs;
using Apocryph.FunctionApp.Model;
using Microsoft.Azure.WebJobs;
using Microsoft.Extensions.Logging;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.FunctionApp
{
    public static class StepSignatureVerifier
    {
        [FunctionName(nameof(StepSignatureVerifier))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("stepsStream")] IAsyncEnumerable<IHashed<IAgentStep>> stepsStream,
            [PerperStream("stepValidatorSetSplitterStream")] IAsyncEnumerable<IHashed<ValidatorSet>> stepValidatorSetSplitterStream,
            [PerperStream("outputStream")] IAsyncCollector<Hash> outputStream,
            ILogger logger)
        {
            await stepsStream.Zip(stepValidatorSetSplitterStream).ForEachAsync(async stepAndValidatorSet =>
            {
                var (step, validatorSet) = stepAndValidatorSet;
                try
                {
                    var commitsSignaturesValid = step.Value.PreviousCommits.All(commit =>
                    {
                        if (commit.Value.For != step.Value.Previous)
                        {
                            return false;
                        }
                        var bytes = IpfsJsonSettings.ObjectToBytes(commit.Value);
                        return commit.Signer.ValidateSignature(bytes, commit.Signature);
                    });

                    if (commitsSignaturesValid)
                    {
                        var signers = step.Value.PreviousCommits.Select(commit => commit.Signer).Distinct();
                        if (validatorSet.Value.IsMoreThanTwoThirds(signers))
                        {
                            await outputStream.AddAsync(step.Value.Previous);
                        }
                    }
                }
                catch (Exception e)
                {
                    logger.LogError(e.ToString());
                }
            }, CancellationToken.None);
        }
    }
}