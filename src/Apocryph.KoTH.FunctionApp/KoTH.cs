using System.Collections.Generic;
using System.Threading.Tasks;
using Apocryph.Consensus;
using Apocryph.Ipfs;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Model;
using Perper.WebJobs.Extensions.Triggers;

namespace Apocryph.KoTH.FunctionApp
{
    public static class KoTH
    {
        [FunctionName("Apocryph-KoTH")]
        public static async Task<(string, IAsyncEnumerable<(Hash<Chain>, Slot?[])>)> Start([PerperTrigger] object? input, IContext context)
        {
            var (minedKeysStream, minedKeysStreamName) = await context.CreateBlankStreamAsync<(Hash<Chain>, Slot)>();

            var resultStream = await context.StreamFunctionAsync<(Hash<Chain>, Slot?[])>("Processor", minedKeysStream);

            return (minedKeysStreamName, resultStream);
        }

        [FunctionName("Processor")]
        public static async IAsyncEnumerable<(Hash<Chain>, Slot?[])> Processor(
            [PerperTrigger] IAsyncEnumerable<(Hash<Chain>, Slot)> minedKeys,
            IState state,
            IHashResolver hashResolver)
        {
            await foreach (var (chain, slot) in minedKeys)
            {
                var chainState = await state.GetValue<KoTHState?>(chain.ToString(), () => default!);
                if (chainState == null)
                {
                    var chainValue = await hashResolver.RetrieveAsync(chain);
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