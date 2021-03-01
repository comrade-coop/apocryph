using System.Collections.Generic;
using System.Linq;
using System.Numerics;
using System.Threading.Tasks;
using Apocryph.Consensus;
using Apocryph.HashRegistry;
using Apocryph.ServiceRegistry;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Model;
using Perper.WebJobs.Extensions.Triggers;

namespace Apocryph.KoTH.FunctionApp
{
    public static class KoTH
    {
        public class KoTHChain
        {
            public Peer?[] Slots { get; }

            public KoTHChain(Peer?[] slots)
            {
                Slots = slots;
            }

            public static BigInteger GetDifficulty(Peer peer)
            {
                var hash = Hash.From(peer);
                return new BigInteger(hash.Bytes.Concat(new byte[] { 0 }).ToArray());
            }

            public bool TryInsert(Peer newPeer)
            {
                var newDifficulty = GetDifficulty(newPeer);
                var slot = (int)(newDifficulty % Slots.Length);

                var currentPeer = Slots[slot];
                if (currentPeer == null || GetDifficulty(currentPeer) < newDifficulty)
                {
                    Slots[slot] = newPeer;
                    return true;
                }

                return false;
            }
        }

        [FunctionName("Apocryph-KoTH")]
        public static async Task Start([PerperTrigger] (IAgent serviceRegistry, HashRegistryProxy hashRegistry) input, IContext context)
        {
            var (minedKeys, minedKeysStream) = await context.CreateBlankStreamAsync<(Hash<Chain>, Peer)>();

            var resultStream = await context.StreamFunctionAsync<(Hash<Chain>, Peer)>("Processor", (minedKeys, input.hashRegistry));

            await input.serviceRegistry.CallActionAsync("Register", (new ServiceLocator("KoTH", "KoTH"), new Service(new Dictionary<string, string>() {
                {"minedKeys", minedKeysStream}
            }, new Dictionary<string, IStream>() {
                {"states", (IStream)resultStream}
            })));
        }

        [FunctionName("Processor")]
        public static async IAsyncEnumerable<(Hash<Chain>, Peer?[])> Processor([PerperTrigger] (IAsyncEnumerable<(Hash<Chain>, Peer)> minedKeys, HashRegistryProxy hashRegistry) input, IState state)
        {
            await foreach (var (chain, peer) in input.minedKeys)
            {
                var chainState = await state.GetValue<KoTHChain?>(chain.ToString(), () => default!);
                if (chainState == null)
                {
                    var chainValue = await input.hashRegistry.RetrieveAsync(chain);
                    chainState = new KoTHChain(new Peer?[chainValue.SlotsCount]);
                }

                if (chainState.TryInsert(peer))
                {
                    await state.SetValue(chain.ToString(), chainState);
                    yield return (chain, chainState.Slots);
                }
            }
        }
    }
}