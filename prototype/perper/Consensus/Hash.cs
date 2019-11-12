using System;
using System.Threading.Tasks;

namespace Apocryph.Consensus
{
    // IPFSHashMessage
    // RevealableHash
    public class Hash<T> where T : new()
    {
        // public byte[] Bytes { get; }

        // GetOriginalValue
        public Task<T> GetValue() {
            throw new Exception();
        }
    }

    public static class Hash
    {
        public static Hash<T> Create<T>(T value) where T : new() {
            throw new Exception();
        }

        public static Hash<T> FromString<T>(string value) where T : new() {
            throw new Exception();
        }
    }
}
