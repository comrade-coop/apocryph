using System.Collections.Generic;
using System.Linq;
using System.Numerics;
using System.Security.Cryptography;
using Apocryph.Core.Consensus.Blocks;

namespace Apocryph.Core.Consensus
{
    public static class RandomWalk
    {
        public static IEnumerable<(BigInteger, Hash)> Run(Hash hash)
        {
            using var sha256Hash = SHA256.Create();
            while (true)
            {
                hash = Hash.From(sha256Hash, hash);
                yield return (new BigInteger(hash.Value.Concat(new byte[] { 0 }).ToArray()), hash);
            }
        }
    }
}