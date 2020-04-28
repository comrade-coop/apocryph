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
    public static class TestValidatorLauncher
    {
        [FunctionName(nameof(TestValidatorLauncher))]
        public static async Task Run([PerperStreamTrigger(RunOnStartup = true)] PerperStreamContext context,
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
            var agentId = "0";

            var _validatorSetsStream = await context.StreamFunctionAsync(nameof(TestDataGenerator), new
            {
                delay = TimeSpan.FromSeconds(10),
                data = genesisValidatorSet
            });

            var validatorSetsStream = await context.StreamFunctionAsync(nameof(IpfsSaver), new
            {
                dataStream = _validatorSetsStream,
                ipfsGateway
            });

            var otherValidatorSetsStream = await context.StreamFunctionAsync("TestValidatorLauncher-Helper", new
            {
                agentId,
                validatorSetsStream,
            });

            await using var validatorLauncherStreams = new AsyncDisposableList();
            foreach (var (privateKey, self) in keys)
            {
                validatorLauncherStreams.Add(
                    await context.StreamActionAsync(nameof(ValidatorLauncher), new
                    {

                        agentId,
                        services = new [] {"Sample", "IpfsInput"},
                        validatorSetsStream = validatorSetsStream,
                        otherValidatorSetsStream,
                        ipfsGateway,
                        privateKey,
                        self
                    }));
            }

            await context.BindOutput(cancellationToken);
        }

        [FunctionName("TestValidatorLauncher-Helper")]
        public static async Task Helper([PerperStreamTrigger] PerperStreamContext context,
            [Perper("agentId")] string agentId,
            [PerperStream("validatorSetsStream")] IAsyncEnumerable<IHashed<ValidatorSet>> validatorSetsStream,
            [PerperStream("outputStream")] IAsyncCollector<Dictionary<string, IHashed<ValidatorSet>>> outputStream)
        {
            await validatorSetsStream.ForEachAsync(async validatorSet =>
            {
                outputStream.AddAsync(new Dictionary<string, IHashed<ValidatorSet>>
                {
                    [agentId] = validatorSet
                });
            }, CancellationToken.None);
        }
    }
}