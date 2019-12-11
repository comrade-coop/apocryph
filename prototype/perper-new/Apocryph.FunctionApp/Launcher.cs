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
        public static async Task Run([PerperTrigger("Launcher")] IPerperStreamContext context)
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

            var commitsStream = ipfsStream; // .Filter(typeof(Commit))
            var votesStream = ipfsStream; // .Filter(typeof(Vote))
            var proposalsStream = ipfsStream; // .Filter(typeof(IAgentStep))

            /* FIXME:
                _committerStream has type (Hash, bool)
                IpfsLoader.hashStream has type Hash
                committerStream has type Hashed<object>
                Runtime.committerStream has type (Hashed<object>, bool) */
            /* FIXME:
                runtimeStream has type (object, bool)
                IpfsSaver.dataStream has type object
                savedRuntimeStream has type Hashed<object>
                Voting.runtimeStream has type (Hashed<object>, bool) */

            var _committerStream = await context.CallStreamFunction("Committer", new
            {
                self,
                commitsStream,
                validatorSet
            });

            var committerStream = await context.CallStreamFunction("IpfsLoader", new
            {
                ipfsGateway,
                hashStream = _committerStream
            });

            var _validatorStream = await context.CallStreamFunction("Validator", new
            {
                commitsStream,
                proposalsStream,
                validatorSet
            });

            var validatorStream = await context.CallStreamFunction("IpfsLoader", new
            {
                ipfsGateway,
                hashStream = _validatorStream
            });

            var runtimeStream = await context.CallStreamFunction("Runtime", new
            {
                validatorStream,
                committerStream
            });

            var proposerStream = await context.CallStreamFunction("Proposer", new
            {
                commitsStream,
                runtimeStream
            });

            var savedRuntimeStream = await context.CallStreamFunction("IpfsSaver", new
            {
                ipfsGateway,
                dataStream = runtimeStream
            });

            var votingStream = await context.CallStreamFunction("Voting", new
            {
                runtimeStream = savedRuntimeStream,
                proposalsStream
            });

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