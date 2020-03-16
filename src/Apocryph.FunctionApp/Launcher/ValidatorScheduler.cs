using System;
using System.Collections.Generic;
using System.Linq;
using System.Security.Cryptography;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Model;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.FunctionApp
{
    public static class ValidatorScheduler
    {
        [FunctionName(nameof(ValidatorScheduler))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("validatorSetsStream")] IAsyncEnumerable<Dictionary<string, IHashed<ValidatorSet>>> validatorSetsStream,
            [Perper("validatorSetsStream")] object[] validatorSetsStreamPassthrough,
            [Perper("ipfsGateway")] string ipfsGateway,
            [Perper("privateKey")] ECParameters privateKey,
            [Perper("self")] ValidatorKey self)
        {
            var runningStreams = new Dictionary<string, IAsyncDisposable>();

            await validatorSetsStream.ForEachAsync(async validatorSets =>
            {
                // FIXME: Instead of restarting when validator set changes, send validator sets as a seperate stream

                var filteredValidatorSets = validatorSets
                    .Where(kv => kv.Value.Value.Weights.ContainsKey(self));

                var toStop = new HashSet<string>(runningStreams.Keys);

                foreach (var kv in filteredValidatorSets)
                {
                    toStop.Remove(kv.Key);

                    if (!runningStreams.ContainsKey(kv.Key))
                    {
                        var agentId = kv.Key;
                        var filterStream = await context.StreamFunctionAsync(nameof(AgentZeroValidatorSetsSplitter), new
                        {
                            validatorSetsStream = validatorSetsStreamPassthrough,
                            agentId,
                        });
                        var initValidatorSetStream = await context.StreamFunctionAsync(nameof(TestDataGenerator), new
                        {
                            delay = TimeSpan.FromSeconds(10),
                            data = kv.Value
                        });
                        var launcher = await context.StreamActionAsync(nameof(ValidatorLauncher), new
                        {
                            agentId = agentId,
                            services = new [] {"Sample", "IpfsInput"},
                            ipfsGateway,
                            validatorSetsStream = new [] {filterStream, initValidatorSetStream},
                            privateKey,
                            self
                        });
                        runningStreams[kv.Key] = new AsyncDisposableList
                        {
                            { filterStream },
                            { initValidatorSetStream },
                            { launcher }
                        };
                    }
                }

                foreach (var kv in toStop)
                {
                    await runningStreams[kv].DisposeAsync();
                    runningStreams.Remove(kv);
                }
            }, CancellationToken.None);
        }
    }
}