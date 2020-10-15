using System;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Core.Consensus.VirtualNodes;

namespace Apocryph.Core.Consensus
{
    public static class Miner
    {
        public static Task RunAsync(Peer self, int proofLength, Assigner assigner, CancellationToken cancellationToken = default)
        {
            // TODO: Maybe run multiple threads in parallel
            return Task.Run(async () =>
            {
                var random = new Random(); // TODO: PRNG should be enough in general, but consider CSRNG just in case
                var buffer = new byte[proofLength];
                try
                {
                    while (!cancellationToken.IsCancellationRequested)
                    {
                        random.NextBytes(buffer);
                        var attempt = buffer.ToArray();

                        assigner.ProcessClaim(self, attempt);

                        await Task.Delay(100);
                    }
                }
                catch (Exception err)
                {
                    Console.WriteLine(err.ToString());
                    throw;
                }
            });
        }
    }
}