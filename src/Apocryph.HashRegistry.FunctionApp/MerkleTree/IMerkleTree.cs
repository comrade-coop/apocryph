using System;
using System.Linq;
using System.Collections.Generic;
using System.Text.Json;
using System.Threading.Tasks;
using Apocryph.HashRegistry.Serialization;

namespace Apocryph.HashRegistry.MerkleTree
{
    public interface IMerkleTree<T>
    {
        IAsyncEnumerable<T> EnumerateItems(HashRegistryProxy proxy);
    }
}