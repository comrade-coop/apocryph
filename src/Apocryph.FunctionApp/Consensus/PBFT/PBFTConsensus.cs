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
    public static class PBFTConsensus
    {
        [FunctionName(nameof(PBFTConsensus))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("agentId")] string agentId,
            [Perper("services")] string[] services,
            [Perper("validatorSetsStream")] object[] validatorSetsStream,
            [Perper("genesisMessage")] (string, object) genesisMessage,
            [Perper("ipfsGateway")] string ipfsGateway,
            [Perper("privateKey")] ECParameters privateKey,
            [Perper("self")] ValidatorKey self,
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

            await using var commandsStream = await context.StreamFunctionAsync(nameof(CommandSplitter), new
            {
                stepsStream = stepOrderVerifierStream
            });

            var commandExecutorStream =  await context.StreamFunctionAsync(nameof(CommandExecutor), new
            {
                commandsStream,
                agentId,
                services,
                genesisMessage,
                ipfsGateway
            });

            // Proposer

            await using var proposerFilterStream = await context.StreamFunctionAsync(nameof(ProposerFilter), new
            {
                self,
                stepsStream = stepOrderVerifierStream,
                currentProposerStream,
            });

            await using var proposerRuntimeStream = await context.StreamFunctionAsync(nameof(Runtime), new
            {
                self,
                agentId,
                inputStream = proposerFilterStream,
            });

            await using var inputProposerStream = await context.StreamFunctionAsync(nameof(ProposerInputProposer), new
            {
                proposerFilterStream,
                commandExecutorStream
            });

            await using var _proposerCommitInjectorStream = await context.StreamFunctionAsync(nameof(ProposerCommitInjector), new
            {
                commitsStream,
                validatorSetsStream,
                stepsStream = new [] {proposerRuntimeStream, inputProposerStream}
            });

            await using var proposerCommitInjectorStream = await context.StreamFunctionAsync(nameof(IpfsSaver), new
            {
                ipfsGateway,
                dataStream = _proposerCommitInjectorStream
            });

            await using var proposerStream = await context.StreamFunctionAsync(nameof(Proposer), new
            {
                proposerCommitInjectorStream
            });

            // Validator

            await using var _proposedStepsStream = await context.StreamFunctionAsync(nameof(ValidatorFilter), new
            {
                currentProposerStream,
                proposalsStream
            });

            await using var proposedStepsStream = await context.StreamFunctionAsync(nameof(IpfsLoader), new
            {
                ipfsGateway,
                hashStream = _proposedStepsStream
            });

            await using var _proposalValidatorSetSplitter = await context.StreamFunctionAsync(nameof(StepValidatorSetSplitter), new
            {
                stepsStream = proposedStepsStream,
            });

            await using var proposalValidatorSetSplitter = await context.StreamFunctionAsync(nameof(IpfsLoader), new
            {
                ipfsGateway,
                hashStream = _proposalValidatorSetSplitter
            });

            await using var proposedStepSignatureVerifierStream = await context.StreamFunctionAsync(nameof(StepSignatureVerifier), new
            {
                stepsStream = proposedStepsStream,
                stepValidatorSetSplitterStream = proposalValidatorSetSplitter
            });

            await using var validatorFilterStream = await context.StreamFunctionAsync(nameof(ValidatorStepOrderVerifier), new
            {
                stepsStream = stepOrderVerifierStream,
                proposedStepsStream = proposedStepSignatureVerifierStream,
                validatorSetsStream
            });

            await using var _validatorRuntimeInputStream = await context.StreamFunctionAsync(nameof(ValidatorRuntimeInput), new
            {
                validatorFilterStream
            });

            await using var validatorRuntimeInputStream = await context.StreamFunctionAsync(nameof(IpfsLoader), new
            {
                ipfsGateway,
                hashStream = _validatorRuntimeInputStream
            });

            await using var validatorRuntimeStream = await context.StreamFunctionAsync(nameof(Runtime), new
            {
                self,
                agentId,
                inputStream = validatorRuntimeInputStream
            });

            await using var _validatorRuntimeOutputStream = await context.StreamFunctionAsync(nameof(ValidatorRuntimeOutput), new
            {
                self,
                validatorFilterStream,
                runtimeStream = validatorRuntimeStream
            });

            await using var validatorRuntimeOutputStream = await context.StreamFunctionAsync(nameof(IpfsSaver), new
            {
                ipfsGateway,
                dataStream = _validatorRuntimeOutputStream
            });

            await using var validatorRuntimeVotingStream = await context.StreamFunctionAsync(nameof(ValidatorRuntimeVoting), new
            {
                validatorRuntimeOutputStream,
                validatorFilterStream,
                genesisStepStream
            });

            await using var inputValidatorStream = await context.StreamFunctionAsync(nameof(ValidatorInputVoting), new
            {
                validatorFilterStream,
                committerStream,
                commandExecutorStream
            });

            // Counting Votes

            await using var voteCounterStream = await context.StreamFunctionAsync(nameof(VoteCounter), new
            {
                validatorSetsStream,
                votesStream
            });

            // Output

            await using var signerStream = await context.StreamFunctionAsync(nameof(Signer), new
            {
                self,
                privateKey,
                dataStream = new[] {proposerStream, validatorRuntimeVotingStream, inputValidatorStream, voteCounterStream}
            });

            await using var ipfsOutputStream = await context.StreamActionAsync(nameof(IpfsOutput), new
            {
                ipfsGateway,
                topic,
                dataStream = signerStream
            });

            await context.BindOutput(cancellationToken);
        }
    }
}