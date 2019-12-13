using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Model;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.FunctionApp
{
    public static class Committer
    {
        private class State
        {
            public Dictionary<Hash, Dictionary<ValidatorKey, ValidatorSignature>> Commits { get; set; }
        }

        [FunctionName("Committer")]
        public static async Task Run([PerperStream("Committer")] IPerperStreamContext context,
            [Perper("validatorSet")] ValidatorSet validatorSet,
            [Perper("commitsStream")] IAsyncEnumerable<Signed<Commit>> commitsStream,
            [Perper("outputStream")] IAsyncCollector<Hash> outputStream)
        {
            var state = context.GetState<State>();

            await commitsStream.ForEachAsync(async commit =>
                {
                    state.Commits[commit.Value.For].Add(commit.Signer, commit.Signature);
                    await context.SaveState();

                    var committed = state.Commits[commit.Value.For].Keys
                        .Select(signer => validatorSet.Weights[signer]).Sum();
                    if (3 * committed > 2 * validatorSet.Total)
                    {
                        await outputStream.AddAsync(commit.Value.For);
                    }
                }, CancellationToken.None);
        }
    }
}