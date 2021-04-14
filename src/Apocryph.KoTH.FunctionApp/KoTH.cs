using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Channels;
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
        public static string PubSubPath = "koth";

        [FunctionName("Apocryph-KoTH")]
        public static async Task<IAsyncEnumerable<(Hash<Chain>, Slot?[])>> Start([PerperTrigger] object? input, IContext context)
        {
            return await context.StreamFunctionAsync<(Hash<Chain>, Slot?[])>("Processor", null);
        }

        [FunctionName("Processor")]
        public static async Task<IAsyncEnumerable<(Hash<Chain>, Slot?[])>> Processor(
            [PerperTrigger] object? input,
            IState state,
            IHashResolver hashResolver,
            IPeerConnector peerConnector,
            CancellationToken cancellationToken)
        {
            var output = Channel.CreateUnbounded<(Hash<Chain>, Slot?[])>();
            var semaphore = new SemaphoreSlim(1, 1); // NOTE: Should use Perper for locking instead
            await peerConnector.ListenPubSub<(Hash<Chain> chain, Slot slot)>(PubSubPath, async (_, message) =>
            {
                await semaphore.WaitAsync();

                var chainState = await state.GetValue<KoTHState?>(message.chain.ToString(), () => default!);
                if (chainState == null)
                {
                    var chainValue = await hashResolver.RetrieveAsync(message.chain);
                    chainState = new KoTHState(new Slot?[chainValue.SlotsCount]);
                }

                if (chainState.TryInsert(message.slot))
                {
                    await state.SetValue(message.chain.ToString(), chainState);
                    await output.Writer.WriteAsync((message.chain, chainState.Slots.ToArray())); // DEBUG: ToArray used due to in-place modifications
                }

                semaphore.Release();
                return true;
            }, cancellationToken);

            cancellationToken.Register(() => output.Writer.Complete()); // DEBUG: Used for testing purposes mainly

            return output.Reader.ReadAllAsync();
        }
    }
}