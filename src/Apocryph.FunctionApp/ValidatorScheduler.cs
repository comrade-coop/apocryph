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
            [PerperStream("validatorSetsStream")] IAsyncEnumerable<Dictionary<string, ValidatorSet>> validatorSetsStream,
            [Perper("ipfsGateway")] string ipfsGateway,
            [Perper("privateKey")] ECParameters privateKey,
            [Perper("self")] ValidatorKey self)
        {
            var runningStreams = new Dictionary<KeyValuePair<string, ValidatorSet>, IAsyncDisposable>();

            await validatorSetsStream.ForEachAsync(async validatorSets =>
            {
                // FIXME: Instead of restarting when validator set changes, send validator sets as a seperate stream

                var filteredValidatorSets = validatorSets
                    .Where(kv => kv.Value.Weights.ContainsKey(self));

                var toStop = new HashSet<KeyValuePair<string, ValidatorSet>>(runningStreams.Keys);

                foreach (var kv in filteredValidatorSets)
                {
                    toStop.Remove(kv);

                    if (!runningStreams.ContainsKey(kv))
                    {
                        var agentId = kv.Key;
                        var validatorSet = kv.Value;
                        runningStreams[kv] = await context.StreamActionAsync(nameof(ValidatorLauncher), new
                        {
                            ipfsGateway,
                            agentId = agentId,
                            validatorSet = validatorSet,
                            privateKey,
                            self
                        });
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