using System.Security.Cryptography;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Model;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.FunctionApp
{
    public static class ValidatorLauncher
    {
        [FunctionName("ValidatorLauncher")]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("agentId")] string agentId,
            [Perper("validatorSet")] ValidatorSet validatorSet,
            [Perper("ipfsGateway")] string ipfsGateway,
            [Perper("privateKey")] ECParameters privateKey,
            [Perper("self")] ValidatorKey self)
        {
            var topic = "apocryph-agent-" + agentId;

            await using var ipfsStream = await context.StreamFunctionAsync("IpfsInput", new
            {
                ipfsGateway,
                topic
            });

            var commitsStream = ipfsStream;
            var votesStream = ipfsStream;
            var proposalsStream = ipfsStream;

            // Proposer (Proposing)

            await using var currentProposerStream = await context.StreamFunctionAsync("CurrentProposer", new
            {
                commitsStream,
                validatorSet
            });

            await using var _committerStream = await context.StreamFunctionAsync("Committer", new
            {
                commitsStream,
                validatorSet
            });

            await using var committerStream = await context.StreamFunctionAsync("IpfsLoader", new
            {
                ipfsGateway,
                hashStream = _committerStream
            });

            await using var commandsStream = await context.StreamFunctionAsync("CommandSplitter", new
            {
                committerStream
            });

            await using var reminderCommandExecutorStream = await context.StreamFunctionAsync("ReminderCommandExecutor", new
            {
                commandsStream
            });

            await using var agentZeroStream = await context.StreamFunctionAsync("IpfsInput", new
            {
                ipfsGateway,
                topic = "apocryph-agent-0"
            });

            await using var _inputVerifierStream = await context.StreamFunctionAsync("StepVerifier", new
            {
                stepsStream = agentZeroStream,
            });

            await using var inputVerifierStream = await context.StreamFunctionAsync("IpfsLoader", new
            {
                ipfsGateway,
                hashStream = _inputVerifierStream
            });

            await using var validatorSetsStream = await context.StreamFunctionAsync("ValidatorSets", new
            {
                inputVerifierStream
            });

            await using var publicationCommandExecutorStream = await context.StreamFunctionAsync("PublicationCommandExecutor", new
            {
                commandsStream,
                validatorSetsStream
            });

            await using var proposerRuntimeStream = await context.StreamFunctionAsync("Runtime", new
            {
                self,
                currentProposerStream,
                committerStream
            });


            await using var inputProposerStream = await context.StreamFunctionAsync("InputProposer", new
            {
                agentInputsStream = new []{reminderCommandExecutorStream, publicationCommandExecutorStream},
                committerStream
            });

            await using var proposerStream = await context.StreamFunctionAsync("Proposer", new
            {
                commitsStream,
                proposerRuntimeStream
            });

            // Validator (Voting)

            await using var _validatorStream = await context.StreamFunctionAsync("Validator", new
            {
                committerStream,
                currentProposerStream,
                proposalsStream,
                validatorSet
            });

            await using var validatorStream = await context.StreamFunctionAsync("IpfsLoader", new
            {
                ipfsGateway,
                hashStream = _validatorStream
            });

            await using var _validatorRuntimeStream = await context.StreamFunctionAsync("Runtime", new
            {
                inputStream = validatorStream
            });

            await using var validatorRuntimeStream = await context.StreamFunctionAsync("IpfsSaver", new
            {
                ipfsGateway,
                dataStream = _validatorRuntimeStream
            });

            await using var votingStream = await context.StreamFunctionAsync("Voting", new
            {
                runtimeStream = validatorRuntimeStream,
                proposalsStream
            });

            // Consensus (Committing)

            await using var consensusStream = await context.StreamFunctionAsync("Consensus", new
            {
                votesStream
            });

            foreach (var stream in new[] {proposerStream, votingStream, consensusStream})
            {
                await using var saverStream = await context.StreamFunctionAsync("IpfsSaver", new
                {
                    ipfsGateway,
                    dataStream = stream
                });

                await using var signerStream = await context.StreamFunctionAsync("Signer", new
                {
                    self,
                    privateKey,
                    dataStream = saverStream
                });

                await context.StreamActionAsync("IpfsOutput", new
                {
                    ipfsGateway,
                    topic,
                    dataStream = signerStream
                });


            }
        }
    }
}