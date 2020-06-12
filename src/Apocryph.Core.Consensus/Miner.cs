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
            using var dsa = ECDsa.Create();
            // TODO: Maybe run multiple threads in parallel
            return Task.Run(() =>
            {
                while (!cancellationToken.IsCancellationRequested)
                {
                    dsa.GenerateKey(PrivateKey.Curve);
                    var privateKey = new PrivateKey(dsa.ExportParameters(true));

                    assigner.AddKey(privateKey.PublicKey, privateKey);
                }
            });
        }
    }
}