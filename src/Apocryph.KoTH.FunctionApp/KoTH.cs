using System.Collections.Generic;
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
        [FunctionName("Apocryph-KoTH")]
        public static async Task Start([PerperTrigger] (IAgent serviceRegistry, HashRegistryProxy hashRegistry) input, IContext context)
        {
            var (minedKeys, minedKeysStream) = await context.CreateBlankStreamAsync<(Hash<Chain>, Slot)>();

            var resultStream = await context.StreamFunctionAsync<(Hash<Chain>, Slot?[])>("Processor", (minedKeys, input.hashRegistry));

            await input.serviceRegistry.CallActionAsync("Register", (new ServiceLocator("KoTH", "KoTH"), new Service(new Dictionary<string, string>() {
                {"minedKeys", minedKeysStream}
            }, new Dictionary<string, IStream>() {
                {"states", (IStream)resultStream}
            })));
        }

        [FunctionName("Processor")]
        public static async IAsyncEnumerable<(Hash<Chain>, Slot?[])> Processor([PerperTrigger] (IAsyncEnumerable<(Hash<Chain>, Slot)> minedKeys, HashRegistryProxy hashRegistry) input, IState state)
        {
            await foreach (var (chain, slot) in input.minedKeys)
            {
                var chainState = await state.GetValue<KoTHState?>(chain.ToString(), () => default!);
                if (chainState == null)
                {
                    var chainValue = await input.hashRegistry.RetrieveAsync(chain);
                    chainState = new KoTHState(new Slot?[chainValue.SlotsCount]);
                }

                if (chainState.TryInsert(slot))
                {
                    await state.SetValue(chain.ToString(), chainState);
                    yield return (chain, chainState.Slots);
                }
            }
        }
    }
}