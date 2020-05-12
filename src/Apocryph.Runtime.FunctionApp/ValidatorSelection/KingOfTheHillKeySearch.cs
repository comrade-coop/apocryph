// NOTE: File is ignored by .csproj file

using System;
using System.Collections.Generic;
using System.Linq;
using System.Security.Cryptography;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Agent;
using Ipfs;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Runtime.FunctionApp.ValidatorSelection
{
    public static class KingOfTheHillKeySearch
    {
        private class State : KingOfTheHillBase.State
        {
            public HashSet<Cid> AgentIds { get; set; } = new HashSet<Cid>();
        }

        [FunctionName(nameof(KingOfTheHillKeySearch))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("claimsStream")] IAsyncEnumerable<ValidatorSlotClaim> claimsStream,
            [PerperStream("agentIdsStream")] IAsyncEnumerable<Cid> agentIdsStream,
            [PerperStream("saltsStream")] IAsyncEnumerable<(int, byte[])> saltsStream,
            [PerperStream("outputStream")] IAsyncCollector<(ECParameters, Cid)> outputStream,
            CancellationToken cancellationToken)
        {
            var state = await context.FetchStateAsync<State>() ?? new State();

            Task.Run(() => // TODO: Run multiple threads in parallel
            {
                while (true)
                {
                    using var dsa = ECDsa.Create(ECCurve.NamedCurves.nistP521);
                    var publicKey = new ValidatorKey { Key = dsa.ExportParameters(false) };
                    var privateKey = dsa.ExportParameters(true);
                    foreach (var agentId in state.AgentIds)
                    {
                        if (state.AddKey(publicKey, agentId))
                        {
                            outputStream.AddAsync((privateKey, agentId));
                        }
                    }
                }
            });

            await Task.WhenAll(
                agentIdsStream.ForEachAsync(async agentId =>
                {
                    state.AgentIds.Add(agentId);
                    await context.UpdateStateAsync(state);
                }, CancellationToken.None),
                KingOfTheHillBase.Run(context, state, claimsStream, saltsStream));
        }
    }
}