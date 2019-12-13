using System.Security.Cryptography;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Model;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.FunctionApp
{
    public static class Launcher
    {
        [FunctionName("Launcher")]
        public static async Task Run([PerperStream("Launcher")] IPerperStreamContext context)
        {
            ValidatorSet validatorSet = new ValidatorSet();

            ECParameters privateKey;
            ValidatorKey self;

            using (var dsa = ECDsa.Create())
            {
                privateKey = dsa.ExportParameters(true);
                self = new ValidatorKey{Key = dsa.ExportParameters(false)};
            }

            var ipfsGateway = "127.0.0.1:5001";
            var topic = "apocryph-test-agent";

            var ipfsStream = await context.CallStreamFunction("IpfsInput", new
            {
                ipfsGateway,
                topic
            });

            var commitsStream = ipfsStream;
            var votesStream = ipfsStream;
            var proposalsStream = ipfsStream;

            // Proposer (Proposing)

            var currentProposerStream = await context.CallStreamFunction("CurrentProposer", new
            {
                commitsStream,
                validatorSet
            });

            var _committerStream = await context.CallStreamFunction("Committer", new
            {
                commitsStream,
                validatorSet
            });

            var committerStream = await context.CallStreamFunction("IpfsLoader", new
            {
                ipfsGateway,
                hashStream = _committerStream
            });

            var proposerRuntimeStream = await context.CallStreamFunction("ProposerRuntime", new
            {
                self,
                currentProposerStream,
                committerStream
            });

            var proposerStream = await context.CallStreamFunction("Proposer", new
            {
                commitsStream,
                proposerRuntimeStream
            });

            // Validator (Voting)

            var _validatorStream = await context.CallStreamFunction("Validator", new
            {
                committerStream,
                currentProposerStream,
                proposalsStream,
                validatorSet
            });

            var validatorStream = await context.CallStreamFunction("IpfsLoader", new
            {
                ipfsGateway,
                hashStream = _validatorStream
            });

            var _validatorRuntimeStream = await context.CallStreamFunction("ProposerRuntime", new
            {
                validatorStream,
                committerStream
            });

            var validatorRuntimeStream = await context.CallStreamFunction("IpfsSaver", new
            {
                ipfsGateway,
                dataStream = _validatorRuntimeStream
            });

            var votingStream = await context.CallStreamFunction("Voting", new
            {
                runtimeStream = validatorRuntimeStream,
                proposalsStream
            });

            // Consensus (Committing)

            var consensusStream = await context.CallStreamFunction("Consensus", new
            {
                votesStream
            });

            foreach (var stream in new[] {proposerStream, votingStream, consensusStream})
            {
                var saverStream = await context.CallStreamFunction("IpfsSaver", new
                {
                    ipfsGateway,
                    dataStream = stream
                });

                var signerStream = await context.CallStreamFunction("Signer", new
                {
                    self,
                    privateKey,
                    dataStream = saverStream
                });

                await context.CallStreamAction("IpfsOutput", new
                {
                    ipfsGateway,
                    topic,
                    dataStream = signerStream
                });
            }
        }
    }
}