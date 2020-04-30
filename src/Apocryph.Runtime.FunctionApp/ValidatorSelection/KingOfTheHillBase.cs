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
            public byte[][] Salts { get; set; } = new byte[0][];
            public ValidatorKey?[] Slots { get; set; } = new ValidatorKey?[0];

            public bool AddKey(ValidatorKey key)
            {
                var slot = (int)(new BigInteger(key.GetPosition()) % positions);
                var salt = Salts[slot];

                if (!(Slots[slot] is ValidatorKey slotKey && new BigInteger(slotKey.GetDifficulty(salt)) < new BigInteger(key.GetDifficulty(salt)))) {
                    Slots[slot] = key;
                    return true;
                }
                return false;
            }
        }

        public static async Task Run<T>(PerperStreamContext context, T state, IAsyncEnumerable<ValidatorKey> seenKeysStream, IAsyncEnumerable<(int, byte[])> saltsStream)
            where T : State
        {
            if (state.Slots.Length != positions)
            {
                state.Slots = new ValidatorKey?[positions];
            }

            if (state.Salts.Length != positions)
            {
                state.Salts = new byte[positions][];
            }

            await Task.WhenAll(
                saltsStream.ForEachAsync(async item =>
                {
                    var (slot, salt) = item;
                    state.Salts[slot] = salt;
                    await context.UpdateStateAsync(state);
                }, CancellationToken.None),

                seenKeysStream.ForEachAsync(async key =>
                {
                    if (state.AddKey(key))
                    {
                        await context.UpdateStateAsync(state);
                    }
                }, CancellationToken.None));
        }
    }
}