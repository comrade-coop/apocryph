using System;
using System.Collections.Generic;
using System.Linq;
using System.Security.Cryptography;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Model;
using Apocryph.FunctionApp.Ipfs;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.FunctionApp.Launcher
{
    public static class TestNodeLauncher
    {
        [FunctionName(nameof(TestNodeLauncher))]
        public static async Task Run([PerperStreamTrigger(/*RunOnStartup = true*/)] PerperStreamContext context,
            CancellationToken cancellationToken)
        {
            var keys = new List<(ECParameters, ValidatorKey)>();
            var genesisValidatorSet = new ValidatorSet();

            for (var i = 0; i < 1; i ++)
            {
                using var dsa = ECDsa.Create(ECCurve.NamedCurves.nistP521);
                var privateKey = dsa.ExportParameters(true);
                var publicKey = new ValidatorKey{Key = dsa.ExportParameters(false)};
                keys.Add((privateKey, publicKey));
                genesisValidatorSet.Weights.Add(publicKey, 10);
            }

            var ipfsGateway = "http://127.0.0.1:5001";

            await using var validatorLauncherStreams = new AsyncDisposableList();
            foreach (var (privateKey, self) in keys)
            {
                validatorLauncherStreams.Add(
                    await context.StreamActionAsync(nameof(NodeLauncher), new
                    {
                        genesisValidatorSet,
                        ipfsGateway,
                        privateKey,
                        self
                    }));
            }

            await context.BindOutput(cancellationToken);
        }
    }
}