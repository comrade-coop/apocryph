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
            [PerperStream("outputStream")] IAsyncCollector<IHashed<IAgentStep>> outputStream,
            ILogger logger)
        {
            await using var stepValidatorSetSplitterEnumerator = stepValidatorSetSplitterStream.GetAsyncEnumerator();
            await stepsStream.ForEachAsync(async step =>
            {
                if (step.Value.Previous == new Hash { Bytes = new byte[]{} })
                {
                    return;
                }

                await stepValidatorSetSplitterEnumerator.MoveNextAsync();

                var validatorSet = stepValidatorSetSplitterEnumerator.Current;
                try
                {
                    var commitsSignaturesValid = step.Value.PreviousCommits.All(commit =>
                    {
                        if (commit.Value.For != step.Value.Previous)
                        {
                            return false;
                        }
                        var bytes = IpfsJsonSettings.ObjectToBytes<object>(commit.Value);
                        return commit.Signer.ValidateSignature(bytes, commit.Signature);
                    });

                    if (commitsSignaturesValid)
                    {
                        var signers = step.Value.PreviousCommits.Select(commit => commit.Signer).Distinct();
                        if (validatorSet.Value.IsMoreThanTwoThirds(signers))
                        {
                            await outputStream.AddAsync(step);
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