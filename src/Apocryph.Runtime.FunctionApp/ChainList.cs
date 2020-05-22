using System;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Runtime.FunctionApp.Consensus.Core;
using Apocryph.Runtime.FunctionApp.Consensus;
using Apocryph.Runtime.FunctionApp.ValidatorSelection;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Runtime.FunctionApp
{
    public class ChainList
    {
        [FunctionName(nameof(ChainList))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("chains")] IAsyncDisposable chains,
            CancellationToken cancellationToken)
        {
            var miner = await context.StreamFunctionAsync(typeof(Miner), new {});
            var gossips = context.DeclareStream(typeof(Peering));
            var ibc = context.DeclareStream(typeof(Peering));

            var factory = await context.StreamFunctionAsync(typeof(ChainFactory), new { chains, miner, gossips, ibc });
            await context.StreamFunctionAsync(ibc, new { factory, filter = typeof(Chain) });
            await context.StreamFunctionAsync(gossips, new { factory, filter = typeof(Assigner) });

            var nodes = new Node[] { new Node { Id = 1 } };
            var chain = await context.StreamFunctionAsync(typeof(Chain), new { localNodes = miner });
        }
    }
}