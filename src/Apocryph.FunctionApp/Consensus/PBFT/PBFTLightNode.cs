using System;
using System.Collections.Generic;
using System.Linq;
using System.Security.Cryptography;
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
    public static class PBFTLightNode
    {
        [FunctionName(nameof(PBFTLightNode))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("agentId")] string agentId,
            [Perper("validatorSetsStream")] object[] validatorSetsStream,
            [Perper("ipfsGateway")] string ipfsGateway,
            CancellationToken cancellationToken)
        {
            var topic = "apocryph-agent-" + agentId;

            await using var ipfsStream = await context.StreamFunctionAsync(nameof(IpfsInput), new
            {
                ipfsGateway,
                topic
            });

            var commitsStream = ipfsStream;
            var votesStream = ipfsStream;
            var proposalsStream = ipfsStream;

            // Initialization

            await using var _genesisStepStream = await context.StreamFunctionAsync(nameof(TestDataGenerator), new
            {
                delay = TimeSpan.FromSeconds(21),
                data = new AgentOutput
                {
                    State = new object(),
                    Commands = new List<ICommand>(),
                    Previous = new Hash { Bytes = new byte[]{} },
                    PreviousValidatorSet = new Hash { Bytes = new byte[]{} },
                    PreviousCommits = new List<ISigned<Commit>>()
                }
            });

            await using var genesisStepStream = await context.StreamFunctionAsync(nameof(IpfsSaver), new
            {
                ipfsGateway,
                dataStream = _genesisStepStream
            });

            // Counting commits

            await using var currentProposerStream = await context.StreamFunctionAsync(nameof(CurrentProposer), new
            {
                commitsStream,
                validatorSetsStream
            });

            await using var _committerStream = await context.StreamFunctionAsync(nameof(CommitCounter), new
            {
                commitsStream,
                validatorSetsStream
            });

            await using var committerStream = await context.StreamFunctionAsync(nameof(IpfsLoader), new
            {
                ipfsGateway,
                hashStream = _committerStream
            });

            // - Sync

            await using var _unverifiedStepsStream = await context.StreamFunctionAsync(nameof(StepHashCollector), new
            {
                inputStream = ipfsStream
            });

            await using var unverifiedStepsStream = await context.StreamFunctionAsync(nameof(IpfsRecursiveLoader), new
            {
                ipfsGateway,
                hashStream = _unverifiedStepsStream
            });

            await using var _stepValidatorSetSplitterStream = await context.StreamFunctionAsync(nameof(StepValidatorSetSplitter), new
            {
                stepsStream = unverifiedStepsStream,
            });

            await using var stepValidatorSetSplitterStream = await context.StreamFunctionAsync(nameof(IpfsLoader), new
            {
                ipfsGateway,
                hashStream = _stepValidatorSetSplitterStream
            });

            await using var stepSignatureVerifierStream = await context.StreamFunctionAsync(nameof(StepSignatureVerifier), new
            {
                stepsStream = unverifiedStepsStream,
                stepValidatorSetSplitterStream
            });

            await using var _verifiedStepsStream = await context.StreamFunctionAsync(nameof(StepVerifiedStepGetter), new
            {
                stepSignatureVerifierStream
            });

            await using var verifiedStepsStream = await context.StreamFunctionAsync(nameof(IpfsLoader), new
            {
                ipfsGateway,
                hashStream = _verifiedStepsStream
            });

            await using var stepOrderVerifierStream = await context.StreamFunctionAsync(nameof(StepOrderVerifier), new
            {
                stepsStream = new []{verifiedStepsStream, committerStream},
                validatorSetsStream
            });

            await context.BindOutput(stepOrderVerifierStream, cancellationToken);
        }
    }
}