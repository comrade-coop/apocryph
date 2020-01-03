using System.Linq;
using System.Security.Cryptography;
using System.Threading;
using System.Threading.Tasks;
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
            [Perper("validatorSet")] ValidatorSet validatorSet,
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

            // Proposer (Proposing)

            await using var currentProposerStream = await context.StreamFunctionAsync(nameof(CurrentProposer), new
            {
                commitsStream,
                validatorSet
            });

            await using var _committerStream = await context.StreamFunctionAsync(nameof(Committer), new
            {
                commitsStream,
                validatorSet
            });

            await using var committerStream = await context.StreamFunctionAsync(nameof(IpfsLoader), new
            {
                ipfsGateway,
                hashStream = _committerStream
            });

            await using var commandsStream = await context.StreamFunctionAsync(nameof(CommandSplitter), new
            {
                committerStream
            });

            await using var reminderCommandExecutorStream = await context.StreamFunctionAsync(nameof(ReminderCommandExecutor), new
            {
                commandsStream
            });

            await using var agentZeroStream = await context.StreamFunctionAsync(nameof(IpfsInput), new
            {
                ipfsGateway,
                topic = "apocryph-agent-0"
            });

            await using var _inputVerifierStream = await context.StreamFunctionAsync(nameof(StepVerifier), new
            {
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
                commandsStream,
                validatorSetsStream
            });

            await using var proposerRuntimeStream = await context.StreamFunctionAsync(nameof(Runtime), new
            {
                self,
                currentProposerStream,
                committerStream
            });


            await using var inputProposerStream = await context.StreamFunctionAsync(nameof(InputProposer), new
            {
                agentInputsStream = new []{reminderCommandExecutorStream, subscriptionCommandExecutorStream},
                committerStream
            });

            await using var proposerStream = await context.StreamFunctionAsync(nameof(Proposer), new
            {
                commitsStream,
                proposerRuntimeStream
            });

            // Validator (Voting)

            await using var validatorFilterStream = await context.StreamFunctionAsync(nameof(ValidatorFilter), new
            {
                committerStream,
                currentProposerStream,
                proposalsStream
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

            await using var _validatorRuntimeStream = await context.StreamFunctionAsync(nameof(Runtime), new
            {
                inputStream = validatorRuntimeInputStream
            });

            await using var validatorRuntimeStream = await context.StreamFunctionAsync(nameof(IpfsSaver), new
            {
                ipfsGateway,
                dataStream = _validatorRuntimeStream
            });

            await using var votingStream = await context.StreamFunctionAsync(nameof(Voting), new
            {
                validatorRuntimeStream,
                validatorFilterStream
            });

            // Consensus (Committing)

            await using var consensusStream = await context.StreamFunctionAsync(nameof(Consensus), new
            {
                votesStream
            });

            await using var outputSaverStream = await context.StreamFunctionAsync(nameof(IpfsSaver), new
            {
                ipfsGateway,
                dataStream = new[] {proposerStream, votingStream, consensusStream}
            });

            await using var signerStream = await context.StreamFunctionAsync(nameof(Signer), new
            {
                self,
                privateKey,
                dataStream = outputSaverStream
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