using System.Collections.Generic;
using System.Linq;

namespace Apocryph.HashRegistry.MerkleTree
{
    public struct MerkleTreeLeaf<T> : IMerkleTree<T>
    {
        public MerkleTreeLeaf(T value)
        {
            Value = value;
        }

        public T Value { get; }

        public IAsyncEnumerable<T> EnumerateItems(HashRegistryProxy proxy)
        {
            return new[] { Value }.ToAsyncEnumerable();
        }
    }
}