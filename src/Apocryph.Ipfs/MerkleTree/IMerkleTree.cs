using System.Collections.Generic;

namespace Apocryph.Ipfs.MerkleTree
{
    public interface IMerkleTree<T>
    {
        IAsyncEnumerable<T> EnumerateItems(IHashResolver resolver);
    }
}