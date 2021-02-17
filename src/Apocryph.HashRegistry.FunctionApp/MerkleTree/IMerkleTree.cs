using System.Collections.Generic;

namespace Apocryph.HashRegistry.MerkleTree
{
    public interface IMerkleTree<T>
    {
        IAsyncEnumerable<T> EnumerateItems(HashRegistryProxy proxy);
    }
}