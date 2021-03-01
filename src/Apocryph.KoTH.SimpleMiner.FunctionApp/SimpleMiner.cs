using System;
using System.Collections.Concurrent;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Consensus;
using Apocryph.HashRegistry;
using Apocryph.ServiceRegistry;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Bindings;
using Perper.WebJobs.Extensions.Model;
using Perper.WebJobs.Extensions.Triggers;

namespace Apocryph.KoTH.SimpleMiner.FunctionApp
{
    using KoTHChain = Apocryph.KoTH.FunctionApp.KoTH.KoTHChain;

    public static class SimpleMiner
    {
        [FunctionName("Apocryph-KoTH-SimpleMiner")]
        public static async Task Start([PerperTrigger] IAgent serviceRegistry, IContext context)
        {
            var kothService = await serviceRegistry.CallFunctionAsync<Service>("Lookup", new ServiceLocator("KoTH", "KoTH"));

            await context.CallActionAsync("Miner", (kothService.Inputs["minedKeys"], kothService.Outputs["states"], 0));
        }

        [FunctionName("Miner")]
        public static async Task Miner([PerperTrigger(ParameterExpression = "{'stream': 0}")] (string _, IAsyncEnumerable<(Hash<Chain>, Peer?[])> kothStates, int peerId) input, [Perper(Stream = "{stream}")] IAsyncCollector<(Hash<Chain>, Peer)> minedKeys, CancellationToken cancellationToken)
        {
            var chains = new ConcurrentDictionary<Hash<Chain>, KoTHChain>();

            var generator = Task.Run(async () =>
            {
                var random = new Random();
                while (!cancellationToken.IsCancellationRequested)
                {
                    var attemptData = new byte[16];
                    random.NextBytes(attemptData);
                    var newPeer = new Peer(input.peerId, attemptData);

                    foreach (var (chain, state) in chains)
                    {
                        if (state.TryInsert(newPeer))
                        {
                            await minedKeys.AddAsync((chain, newPeer));
                        }
                    }

                    await Task.Delay(10);
                }
            });

            await foreach (var (chain, peers) in input.kothStates.WithCancellation(cancellationToken))
            {
                chains[chain] = new KoTHChain(peers);
            }

            await generator;
        }
    }
}