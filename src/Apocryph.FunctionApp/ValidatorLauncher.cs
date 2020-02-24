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
    public static class ValidatorLauncher
    {
        [FunctionName(nameof(ValidatorLauncher))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("agentId")] string agentId,
            [Perper("validatorSetsStream")] object[] validatorSetsStream,
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

            await using var initMessageStream = await context.StreamFunctionAsync(nameof(TestDataGenerator), new
            {
                delay = TimeSpan.FromSeconds(10),
                data = ("", (object)new InitMessage())
            });

            await using var _genesisStream = await context.StreamFunctionAsync(nameof(TestDataGenerator), new
            {
                delay = TimeSpan.FromSeconds(15),
                data = new AgentOutput
                {
                    State = new object(),
                    Commands = new List<ICommand>(),
                    Previous = new Hash { Bytes = new byte[]{} },
                    PreviousValidatorSet = new Hash { Bytes = new byte[]{} },
                    PreviousCommits = new List<ISigned<Commit>>()
                }
            });

            await using var genesisStream = await context.StreamFunctionAsync(nameof(IpfsSaver), new
            {
                ipfsGateway,
                dataStream = _genesisStream
            });

            // Committer (Executing)

            await using var currentProposerStream = await context.StreamFunctionAsync(nameof(CurrentProposer), new
            {
                commitsStream,
                validatorSetsStream
            });

            await using var _committerStream = await context.StreamFunctionAsync(nameof(Committer), new
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

            await using var _signatureVerifierStream = await context.StreamFunctionAsync(nameof(StepSignatureVerifier), new
            {
                stepsStream = unverifiedStepsStream,
                stepValidatorSetSplitterStream
            });

            await using var signatureVerifierStream = await context.StreamFunctionAsync(nameof(IpfsLoader), new
            {
                ipfsGateway,
                hashStream = _signatureVerifierStream
            });

            await using var verifiedStepStream = await context.StreamFunctionAsync(nameof(StepOrderVerifier), new
            {
                stepsStream = new []{genesisStream, signatureVerifierStream, committerStream},
                validatorSetsStream
            });

            await using var commandsStream = await context.StreamFunctionAsync(nameof(CommandSplitter), new
            {
                stepsStream = verifiedStepStream
            });

            await using var reminderCommandExecutorStream = await context.StreamFunctionAsync(nameof(ReminderCommandExecutor), new
            {
                commandsStream
            });

            /* await using var agentZeroStream = await context.StreamFunctionAsync(nameof(IpfsInput), new
            {
                ipfsGateway,
                topic = 0 //"apocryph-agent-0"
            });

            await using var _inputVerifierStream = await context.StreamFunctionAsync(nameof(StepSignatureVerifier), new
            {
                validatorSetsStream, // TODO: Should give agent 0's validator set instead !!
                stepsStream = agentZeroStream,
            });

            await using var inputVerifierStream = await context.StreamFunctionAsync(nameof(IpfsLoader), new
            {
                ipfsGateway,
                hashStream = _inputVerifierStream
            });

            await using var validatorSetsStream = await context.StreamFunctionAsync(nameof(ValidatorSets), new
            {
                inputVerifierStream
            });

            await using var subscriptionCommandExecutorStream = await context.StreamFunctionAsync(nameof(SubscriptionCommandExecutor), new
            {
                ipfsGateway,
                commandsStream,
                validatorSetsStream
            }); */

            var commandExecutorStream = new []
            {
                reminderCommandExecutorStream,
                initMessageStream,
                // subscriptionCommandExecutorStream
            };

            // Proposer (Proposing)

            await using var proposerFilterStream = await context.StreamFunctionAsync(nameof(ProposerFilter), new
            {
                self,
                stepsStream = verifiedStepStream,
                currentProposerStream,
            });

            await using var proposerRuntimeStream = await context.StreamFunctionAsync(nameof(Runtime), new
            {
                self,
                inputStream = proposerFilterStream,
            });

            await using var inputProposerStream = await context.StreamFunctionAsync(nameof(InputProposer), new
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

            // Validator (Voting)

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

            await using var _proposalSignatureVerifierStream = await context.StreamFunctionAsync(nameof(StepSignatureVerifier), new
            {
                stepsStream = proposedStepsStream,
                stepValidatorSetSplitterStream = proposalValidatorSetSplitter
            });

            await using var proposalSignatureVerifierStream = await context.StreamFunctionAsync(nameof(IpfsLoader), new
            {
                ipfsGateway,
                hashStream = _proposalSignatureVerifierStream
            });

            await using var validatorFilterStream = await context.StreamFunctionAsync(nameof(ProposedStepOrderVerifier), new
            {
                stepsStream = verifiedStepStream,
                proposedStepsStream = proposalSignatureVerifierStream,
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

            await using var votingStream = await context.StreamFunctionAsync(nameof(Voting), new
            {
                validatorRuntimeOutputStream,
                validatorFilterStream
            });

            await using var inputValidatorStream = await context.StreamFunctionAsync(nameof(InputValidator), new
            {
                validatorFilterStream,
                committerStream,
                commandExecutorStream
            });

            // Consensus (Committing)

            await using var consensusStream = await context.StreamFunctionAsync(nameof(Consensus), new
            {
                validatorSetsStream,
                votesStream
            });

            // Output

            await using var signerStream = await context.StreamFunctionAsync(nameof(Signer), new
            {
                self,
                privateKey,
                dataStream = new[] {proposerStream, votingStream, inputValidatorStream, consensusStream}
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