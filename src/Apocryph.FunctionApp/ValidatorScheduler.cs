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
        [FunctionName("ValidatorScheduler")]
        public static async Task Run([PerperStreamTrigger("ValidatorScheduler")] IPerperStreamContext context,
            [PerperStreamTrigger("validatorSet")] IAsyncEnumerable<Dictionary<string, ValidatorSet>> validatorSetsStream,
            [PerperStreamTrigger("ipfsGateway")] string ipfsGateway,
            [PerperStreamTrigger("privateKey")] string privateKey,
            [PerperStreamTrigger("self")] ValidatorKey self)
        {
            var runningStreams = new Dictionary<KeyValuePair<string, ValidatorSet>, IAsyncDisposable>();

            await validatorSetsStream.ForEachAsync(async validatorSets =>
            {
                // FIXME: Instead of restarting when validator set changes, send validator sets as a seperate stream
                var toStop = new HashSet<KeyValuePair<string, ValidatorSet>>(runningStreams.Keys);
                toStop.ExceptWith(validatorSets);

                foreach (var kv in toStop)
                {
                    await runningStreams[kv].DisposeAsync();
                    runningStreams.Remove(kv);
                }

                foreach (var kv in validatorSets)
                {
                    if (!runningStreams.ContainsKey(kv))
                    {
                        var agentId = kv.Key;
                        var validatorSet = kv.Value;
                        runningStreams[kv] = await context.StreamActionAsync("ValidatorLauncher", new
                        {
                            ipfsGateway,
                            agentId = agentId,
                            validatorSet = validatorSet,
                            privateKey,
                            self
                        }).ContinueWith(x => (IAsyncDisposable)null); // TODO: remove ContinueWith when perper starts returning IAsyncDisposable-s
                    }
                }
            }, CancellationToken.None);
        }
    }
}