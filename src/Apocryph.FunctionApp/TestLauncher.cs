using System.Security.Cryptography;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Model;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.FunctionApp
{
    public static class TestLauncher
    {
        [FunctionName(nameof(TestLauncher))]
        public static async Task Run([PerperStreamTrigger(RunOnStartup = true)] PerperStreamContext context,
            CancellationToken cancellationToken)
        {
            ECParameters privateKey;
            ValidatorKey self;

            using (var dsa = ECDsa.Create(ECCurve.NamedCurves.nistP521))
            {
                privateKey = dsa.ExportParameters(true);
                self = new ValidatorKey{Key = dsa.ExportParameters(false)};
            }

            var ipfsGateway = "http://127.0.0.1:5001";
            var validatorSet = new ValidatorSet();
            validatorSet.Weights.Add(self, 10);

            await using var validatorLauncherStream = await context.StreamActionAsync(nameof(ValidatorLauncher),
                new
                {
                    agentId = "0",
                    validatorSet,
                    ipfsGateway,
                    privateKey, self
                });

            await context.BindOutput(cancellationToken);
        }
    }
}