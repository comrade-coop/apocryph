namespace Apocryph.Ipfs.MerkleTree
{
    public interface IMerkleTreeBuilder<T>
    {
        IMerkleTree<T> Value { get; }
        Hash<IMerkleTree<T>> Hash { get; }
        MerkleTreeNodeBuilder<T>? Parent { get; set; }
        int IndexInParent { get; set; }
    }
}