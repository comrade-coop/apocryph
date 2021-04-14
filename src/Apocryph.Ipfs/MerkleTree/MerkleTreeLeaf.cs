using System.Collections.Generic;
using System.Linq;

namespace Apocryph.Ipfs.MerkleTree
{
    public struct MerkleTreeLeaf<T> : IMerkleTree<T>
    {
        public MerkleTreeLeaf(T value)
        {
            Value = value;
        }

        public T Value { get; private set; }

        public IAsyncEnumerable<T> EnumerateItems(IHashResolver resolver)
        {
            return new[] { Value }.ToAsyncEnumerable();
        }
    }
}