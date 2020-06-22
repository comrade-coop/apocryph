using System.Collections.Generic;
using System.Linq;
using System.Numerics;
using System.Security.Cryptography;

namespace Apocryph.Core.Consensus
{
    public static class RandomWalk
    {
        public static IEnumerable<(BigInteger, byte[])> Run(byte[] hash)
        {
            using var sha256Hash = SHA256.Create();
            while (true)
            {
                hash = sha256Hash.ComputeHash(hash);
                yield return (new BigInteger(hash.Concat(new byte[] { 0 }).ToArray()), hash);
            }
        }
    }
}