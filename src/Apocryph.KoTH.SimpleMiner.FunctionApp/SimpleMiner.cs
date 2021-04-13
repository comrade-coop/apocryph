using System;
using System.Collections.Concurrent;
using System.Collections.Generic;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Consensus;
using Apocryph.Ipfs;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Triggers;

namespace Apocryph.KoTH.SimpleMiner.FunctionApp
{
    public static class SimpleMiner
    {
        public static string PubSubPath = Apocryph.KoTH.FunctionApp.KoTH.PubSubPath;

        [FunctionName("Apocryph-KoTH-SimpleMiner")]
        public static async Task Miner(
            [PerperTrigger] IAsyncEnumerable<(Hash<Chain>, Slot?[])> kothStates,
            IPeerConnector peerConnector,
            CancellationToken cancellationToken)
        {
            var self = await peerConnector.Self;
            var chains = new ConcurrentDictionary<Hash<Chain>, KoTHState>();

            var generator = Task.Run(async () =>
            {
                var random = new Random();
                while (!cancellationToken.IsCancellationRequested)
                {
                    var attemptData = new byte[16];
                    random.NextBytes(attemptData);
                    var newSlot = new Slot(self, attemptData);

                    foreach (var (chain, state) in chains)
                    {
                        if (state.TryInsert(newSlot))
                        {
                            await peerConnector.SendPubSub(PubSubPath, (chain, newSlot), cancellationToken);
                        }
                    }

                    await Task.Delay(5); // DEBUG: Try not to hog a full CPU core while testing
                }
            });

            await foreach (var (chain, peers) in kothStates)
            {
                chains[chain] = new KoTHState(peers);
            }

            await generator;
        }
    }
}