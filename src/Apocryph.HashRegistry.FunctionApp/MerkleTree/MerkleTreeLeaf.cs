using System;
using System.Linq;
using System.Collections.Generic;
using System.Text.Json;
using System.Threading.Tasks;
using Apocryph.HashRegistry.Serialization;

namespace Apocryph.HashRegistry.MerkleTree
{
    public struct MerkleTreeLeaf<T> : IMerkleTree<T>
    {
        public MerkleTreeLeaf(T value)
        {
            Value = value;
        }

        public T Value { get; }

        public IAsyncEnumerable<T> EnumerateItems(HashRegistryProxy proxy) {
            return new[] { Value }.ToAsyncEnumerable();
        }
    }
}