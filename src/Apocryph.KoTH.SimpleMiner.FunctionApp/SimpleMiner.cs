using System;
using System.Collections.Concurrent;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Consensus;
using Apocryph.Ipfs;
using Apocryph.ServiceRegistry;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Bindings;
using Perper.WebJobs.Extensions.Model;
using Perper.WebJobs.Extensions.Triggers;

namespace Apocryph.KoTH.SimpleMiner.FunctionApp
{
    public static class SimpleMiner
    {
        [FunctionName("Apocryph-KoTH-SimpleMiner")]
        public static async Task Start([PerperTrigger] (IAgent serviceRegistry, Peer self) input, IContext context)
        {
            var kothService = await input.serviceRegistry.CallFunctionAsync<Service>("Lookup", new ServiceLocator("KoTH", "KoTH"));

            await context.CallActionAsync("Miner", (kothService.Inputs["minedKeys"], kothService.Outputs["states"], input.self));
        }

        [FunctionName("Miner")]
        public static async Task Miner([PerperTrigger(ParameterExpression = "{'stream': 0}")] (string _, IAsyncEnumerable<(Hash<Chain>, Slot?[])> kothStates, Peer self) input, [Perper(Stream = "{stream}")] IAsyncCollector<(Hash<Chain>, Slot)> minedKeys, CancellationToken cancellationToken)
        {
            var chains = new ConcurrentDictionary<Hash<Chain>, KoTHState>();

            var generator = Task.Run(async () =>
            {
                var random = new Random();
                while (!cancellationToken.IsCancellationRequested)
                {
                    var attemptData = new byte[16];
                    random.NextBytes(attemptData);
                    var newPeer = new Slot(input.self, attemptData);

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
                chains[chain] = new KoTHState(peers);
            }

            await generator;
        }
    }
}