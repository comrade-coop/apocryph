using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Core.Consensus.VirtualNodes;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Runtime.FunctionApp
{
    public class Chain
    {
        [FunctionName(nameof(Chain))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("miner")] IAsyncEnumerable<(PrivateKey, string)> miner,
            [Perper("gossips")] IAsyncDisposable gossips,
            [Perper("queries")] IAsyncDisposable queries,
            [PerperStream("output")] IAsyncCollector<IAsyncDisposable> output,
            CancellationToken cancellationToken)
        {
            await foreach (var (privateKey, chain) in miner.WithCancellation(cancellationToken))
            {
                IAsyncDisposable salts = default!;

                var assigner = await context.StreamFunctionAsync(typeof(Assigner), new {chain, privateKey, gossips, salts});
                var proposer = await context.StreamFunctionAsync(typeof(Proposer), new {assigner, queries});
                var validator = await context.StreamFunctionAsync(typeof(Validator), new {assigner, queries});
                var committer = await context.StreamFunctionAsync(typeof(Committer),
                    new {assigner, gossips, proposer, validator});

                await Task.WhenAll(new[] {assigner, proposer, validator, committer}.Select(
                    stream => output.AddAsync(stream, cancellationToken)));
            }
        }
    }
}