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
        private class State
        {
            public ValidatorSet ValidatorSet { get; set; }
        }

        [FunctionName(nameof(StepSignatureVerifier))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("stepsStream")] IAsyncEnumerable<IHashed<IAgentStep>> stepsStream,
            [PerperStream("validatorSetStream")] IAsyncEnumerable<ValidatorSet> validatorSetStream,
            [PerperStream("outputStream")] IAsyncCollector<Hash> outputStream,
            ILogger logger)
        {

            var state = await context.FetchStateAsync<State>() ?? new State();

            await Task.WhenAll(
                validatorSetStream.ForEachAsync(async validatorSet =>
                {
                    try
                    {
                        state.ValidatorSet = validatorSet;
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
                            var committed = step.Value.PreviousCommits
                                .Select(commit => commit.Signer).Distinct()
                                .Select(signer => state.ValidatorSet.Weights[signer]).Sum();
                            if (3 * committed > 2 * state.ValidatorSet.Total)
                            {
                                await outputStream.AddAsync(step.Value.Previous);
                            }
                        }
                    }
                    catch (Exception e)
                    {
                        logger.LogError(e.ToString());
                    }
                }, CancellationToken.None));
        }
    }
}