using System;
using System.Collections.Concurrent;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Consensus;
using Apocryph.Ipfs;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Bindings;
using Perper.WebJobs.Extensions.Triggers;

namespace Apocryph.KoTH.SimpleMiner.FunctionApp
{
    public static class SimpleMiner
    {
        [FunctionName("Apocryph-KoTH-SimpleMiner")]
        public static async Task Miner(
            [PerperTrigger(ParameterExpression = "{'stream': 0}")] (string _, IAsyncEnumerable<(Hash<Chain>, Slot?[])> kothStates) input,
            [Perper(Stream = "{stream}")] IAsyncCollector<(Hash<Chain>, Slot)> minedKeys,
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
                    var newPeer = new Slot(self, attemptData);

                    foreach (var (chain, state) in chains)
                    {
                        if (state.TryInsert(newPeer))
                        {
                            await minedKeys.AddAsync((chain, newPeer));
                        }
                    }

                    // await Task.Delay(10); // DEBUG: Try not to hog a full CPU core while testing
                }
            });

            await foreach (var (chain, peers) in input.kothStates.WithCancellation(cancellationToken))
            {
                chains[chain] = new KoTHState(peers);
            }

            await generator;
        }
    }
}