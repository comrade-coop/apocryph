// NOTE: File is ignored by .csproj file

using System;
using System.Collections.Generic;
using System.Linq;
using System.Security.Cryptography;
using System.Threading;
using System.Threading.Tasks;
using System.Numerics;
using Apocryph.Agent;
using Apocryph.Runtime.FunctionApp.Ipfs;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Runtime.FunctionApp.ValidatorSelection
{
    public static class KingOfTheHillBase
    {
        const int positions = 100; // 1000 * 1000

        public class State
        {
            public ValidatorKey?[] Slots { get; set; }

            public bool AddKey(ValidatorKey key)
            {
                var slot = (int)(new BigInteger(key.GetPosition()) % positions);

                if (!(Slots[slot] is ValidatorKey slotKey && slotKey.CompareTo(key) < 0)) {
                    Slots[slot] = key;
                    return true;
                }
                return false;
            }
        }

        public static async Task Run<T>(PerperStreamContext context, T state, IAsyncEnumerable<ValidatorKey> seenKeysStream)
            where T : State
        {
            if (state.Slots == null || state.Slots.Length != positions)
            {
                state.Slots = new ValidatorKey?[positions];
            }

            await seenKeysStream.ForEachAsync(async key =>
            {
                if (state.AddKey(key))
                {
                    await context.UpdateStateAsync(state);
                }
            }, CancellationToken.None);
        }
    }
}