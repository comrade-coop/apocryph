using System;

namespace Apocryph.Ipfs.MerkleTree
{
    public class MerkleTreeNodeBuilder<T> : IMerkleTreeBuilder<T>
    {
        public MerkleTreeNode<T> Value { get; private set; }
        IMerkleTree<T> IMerkleTreeBuilder<T>.Value { get => Value; }
        public Hash<IMerkleTree<T>> Hash { get; private set; }
        public MerkleTreeNodeBuilder<T>? Parent { get; set; }
        public int IndexInParent { get; set; }

        public MerkleTreeNodeBuilder(MerkleTreeNode<T> value, Hash<IMerkleTree<T>> hash)
        {
            Value = value;
            Hash = hash;
        }

        public (Hash<IMerkleTree<T>>[] left, Hash<IMerkleTree<T>>[] right) GetLevelProof(int childIndex)
        {
            var leftChildren = new Hash<IMerkleTree<T>>[childIndex];
            Array.Copy(Value.Children, 0, leftChildren, 0, leftChildren.Length);

            var rightChildren = new Hash<IMerkleTree<T>>[Value.Children.Length - childIndex - 1];
            Array.Copy(Value.Children, childIndex + 1, rightChildren, 0, rightChildren.Length);

            return (leftChildren, rightChildren);
        }
    }
}