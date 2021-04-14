using System.Collections.Generic;

namespace Apocryph.Ipfs.MerkleTree
{
    public struct MerkleTreeNode<T> : IMerkleTree<T>
    {
        public MerkleTreeNode(Hash<IMerkleTree<T>>[] children)
        {
            Children = children;
        }

        public Hash<IMerkleTree<T>>[] Children { get; private set; }

        public async IAsyncEnumerable<T> EnumerateItems(IHashResolver resolver)
        {
            foreach (var child in Children)
            {
                var subtree = await resolver.RetrieveAsync(child);
                await foreach (var item in subtree.EnumerateItems(resolver))
                {
                    yield return item;
                }
            }
        }
    }
}