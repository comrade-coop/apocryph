using System.Collections.Generic;

namespace Apocryph.Ipfs.MerkleTree
{
    public class MerkleTreeLeafBuilder<T> : IMerkleTreeBuilder<T>
    {
        public MerkleTreeLeaf<T> Value { get; private set; }
        IMerkleTree<T> IMerkleTreeBuilder<T>.Value { get => Value; }
        public Hash<IMerkleTree<T>> Hash { get; private set; }
        public MerkleTreeNodeBuilder<T>? Parent { get; set; }
        public int IndexInParent { get; set; }

        public MerkleTreeLeafBuilder(MerkleTreeLeaf<T> value, Hash<IMerkleTree<T>> hash)
        {
            Value = value;
            Hash = hash;
        }

        public MerkleTreeProof<T> GetProof()
        {
            var levels = new List<(Hash<IMerkleTree<T>>[], Hash<IMerkleTree<T>>[])>();

            IMerkleTreeBuilder<T>? current = this;
            while (current.Parent != null)
            {
                levels.Add(current.Parent.GetLevelProof(current.IndexInParent));
                current = current.Parent;
            }

            return new MerkleTreeProof<T>(Value.Value, levels.ToArray());
        }

        public IMerkleTree<T> GetRoot()
        {
            IMerkleTreeBuilder<T> result = this;
            while (result.Parent != null)
            {
                result = result.Parent;
            }
            return result.Value;
        }
    }
}