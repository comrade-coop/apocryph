using System;
using System.Linq;
using System.Collections.Generic;
using System.Text.Json;
using System.Threading.Tasks;
using Apocryph.HashRegistry.Serialization;

namespace Apocryph.HashRegistry.MerkleTree
{
    public interface IMerkleTreeBuilder<T>
    {
        IMerkleTree<T> Value { get; }
        Hash<IMerkleTree<T>> Hash { get; }
        MerkleTreeNodeBuilder<T>? Parent { get; set; }
        int IndexInParent { get; set; }
    }
}