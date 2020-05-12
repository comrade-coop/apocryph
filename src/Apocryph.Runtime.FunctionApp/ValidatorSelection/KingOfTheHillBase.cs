// NOTE: File is ignored by .csproj file

using System;
using System.Collections.Generic;
using System.Linq;
using System.Security.Cryptography;
using System.Threading;
using System.Threading.Tasks;
using System.Numerics;
using Apocryph.Agent;
using Ipfs;
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
            public Dictionary<Cid, ValidatorKey?[]> Slots { get; set; } = new Dictionary<Cid, ValidatorKey?[]>();

            public bool AddKey(ValidatorKey key, Cid agentId)
            {
                if (!Slots.ContainsKey(agentId) || Slots[agentId].Length != positions)
                {
                    Slots[agentId] = new ValidatorKey?[positions];
                }

                var agentSlots = Slots[agentId];
                var slot = (int)(new BigInteger(key.GetPosition()) % positions);
                var salt = Salts[slot];

                if (!(agentSlots[slot] is ValidatorKey slotKey
                    && new BigInteger(key.GetDifficulty(agentId, salt)) > new BigInteger(slotKey.GetDifficulty(agentId, salt))))
                {
                    agentSlots[slot] = key;
                    return true;
                }
                return false;
            }
        }

        public static async Task Run<T>(PerperStreamContext context, T state,
            IAsyncEnumerable<ValidatorSlotClaim> claimsStream,
            IAsyncEnumerable<(int, byte[])> saltsStream)
            where T : State
        {
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

                claimsStream.ForEachAsync(async claim =>
                {
                    if (state.AddKey(claim.Key, claim.AgentId))
                    {
                        await context.UpdateStateAsync(state);
                    }
                }, CancellationToken.None));
        }
    }
}