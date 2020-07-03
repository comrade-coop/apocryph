using System;
using System.Security.Cryptography;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Core.Consensus.VirtualNodes;

namespace Apocryph.Core.Consensus
{
    public static class Miner
    {
        public static Task RunAsync(Assigner assigner, CancellationToken cancellationToken = default)
        {
            // TODO: Maybe run multiple threads in parallel
            return Task.Run(async () =>
            {
                using var dsa = ECDsa.Create();
                try
                {
                    while (!cancellationToken.IsCancellationRequested)
                    {
                        dsa.GenerateKey(PrivateKey.Curve);
                        var privateKey = new PrivateKey(dsa.ExportParameters(true));

                        assigner.AddKey(privateKey.PublicKey, privateKey);

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