using System;
using System.Linq;
using System.Collections.Generic;
using System.Text.Json;
using System.Threading.Tasks;
using Apocryph.HashRegistry.Serialization;

namespace Apocryph.HashRegistry.MerkleTree
{
    public class MerkleTreeLeafBuilder<T> : IMerkleTreeBuilder<T>
    {
        public MerkleTreeLeaf<T> Value { get; }
        IMerkleTree<T> IMerkleTreeBuilder<T>.Value { get => Value; }
        public Hash<IMerkleTree<T>> Hash { get; }
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