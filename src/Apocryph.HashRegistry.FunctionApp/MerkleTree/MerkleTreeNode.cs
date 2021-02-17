using System.Collections.Generic;

namespace Apocryph.HashRegistry.MerkleTree
{
    public struct MerkleTreeNode<T> : IMerkleTree<T>
    {
        public MerkleTreeNode(Hash<IMerkleTree<T>>[] children)
        {
            Children = children;
        }

        public Hash<IMerkleTree<T>>[] Children { get; }

        public async IAsyncEnumerable<T> EnumerateItems(HashRegistryProxy proxy)
        {
            foreach (var child in Children)
            {
                var subtree = await proxy.RetrieveAsync(child);
                await foreach (var item in subtree.EnumerateItems(proxy))
                {
                    yield return item;
                }
            }
        }
    }
}