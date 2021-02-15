using System;
using System.Linq;
using System.Collections.Generic;
using System.Text.Json;
using System.Threading.Tasks;
using Apocryph.HashRegistry.Serialization;

namespace Apocryph.HashRegistry.MerkleTree
{
    public static class MerkleTreeBuilder
    {
        public static async Task<MerkleTreeLeafBuilder<T>[]> CreateFromValues<T>(HashRegistryProxy proxy, IEnumerable<T> values, int maxChildren)
        {
            var result = new List<MerkleTreeLeafBuilder<T>>();
            foreach (var value in values)
            {
                var leaf = new MerkleTreeLeaf<T>(value);
                var hash = await proxy.StoreAsync<IMerkleTree<T>>(leaf);
                var builder = new MerkleTreeLeafBuilder<T>(leaf, hash);
                result.Add(builder);
            }

            var currentLayer = new Queue<IMerkleTreeBuilder<T>>(result);

            while (currentLayer.Count > 1)
            {
                var previousLayer = currentLayer;
                currentLayer = new Queue<IMerkleTreeBuilder<T>>((previousLayer.Count - 1) / maxChildren + 1);

                while (previousLayer.Count > 0)
                {
                    if (previousLayer.Count == 1)
                    {
                        currentLayer.Enqueue(previousLayer.Dequeue()); // Optimization; we never want single-element nodes
                        break;
                    }

                    var children = new IMerkleTreeBuilder<T>[Math.Min(previousLayer.Count, maxChildren)];
                    var hashes = new Hash<IMerkleTree<T>>[children.Length];
                    for (var i = 0; i < children.Length; i++)
                    {
                        children[i] = previousLayer.Dequeue();
                        hashes[i] = children[i].Hash;
                    }
                    var node = new MerkleTreeNode<T>(hashes);
                    var hash = await proxy.StoreAsync<IMerkleTree<T>>(node);
                    var builder = new MerkleTreeNodeBuilder<T>(node, hash);
                    for (var i = 0; i < children.Length; i++)
                    {
                        children[i].Parent = builder;
                        children[i].IndexInParent = i;
                    }
                    currentLayer.Enqueue(builder);
                }
            }

            return result.ToArray();
        }
    }
}