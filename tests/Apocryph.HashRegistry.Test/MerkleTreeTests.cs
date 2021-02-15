using System.Collections.Generic;
using System.Linq;
using Perper.WebJobs.Extensions.Fake;
using Apocryph.HashRegistry.MerkleTree;
using Xunit;

namespace Apocryph.HashRegistry.Test
{
    using HashRegistry = Apocryph.HashRegistry.FunctionApp.HashRegistry;

    public class MerkleTreeTests
    {
        private HashRegistryProxy GetHashRegistryProxy()
        {
            var registry = new HashRegistry(new FakeState());

            var agent = new FakeAgent();
            agent.RegisterFunction("Store", (byte[] data) => registry.Store(data, default));
            agent.RegisterFunction("Retrieve", (Hash hash) => registry.Retrieve(hash, default));

            return new HashRegistryProxy(agent);
        }

        class Example
        {
            public int Number { get; set; }

            public override bool Equals(object? obj) => obj is Example other && Number == other.Number;
            public override int GetHashCode() => Number.GetHashCode();
            public override string ToString() => $"Example({Number})";
        }

        [Theory]
        [InlineData(3, 4)]
        public async void Enumerate_ReturnsItems_InOrder(int groups, int elements)
        {
            var proxy = GetHashRegistryProxy();

            var rootHashes = new Hash<IMerkleTree<Example>>[groups];
            var expected = new List<Example>();
            for (var i = 0; i < groups; i ++)
            {
                var groupHashes = new Hash<IMerkleTree<Example>>[elements];
                for (var j = 0; j < elements; j++)
                {
                    var example = new Example { Number = i * elements + j };
                    expected.Add(example);
                    var leaf = new MerkleTreeLeaf<Example>(example);
                    groupHashes[j] = await proxy.StoreAsync<IMerkleTree<Example>>(leaf, default);
                }
                var group = new MerkleTreeNode<Example>(groupHashes);
                rootHashes[i] = await proxy.StoreAsync<IMerkleTree<Example>>(group, default);
            }

            var root = new MerkleTreeNode<Example>(rootHashes);

            var result = await root.EnumerateItems(proxy).ToArrayAsync();

            Assert.Equal(expected.ToArray(), result, new ArrayComparer<Example>());
        }

        [Theory]
        [InlineData(8, 2)]
        [InlineData(15, 2)]
        [InlineData(5, 3)]
        [InlineData(7, 3)]
        public async void Builder_ReturnsItems_InOrder(int elements, int maxChildrenCount)
        {
            var proxy = GetHashRegistryProxy();

            var expected = Enumerable.Range(0, elements).Select(x => new Example { Number = x }).ToArray();
            var proofs = await MerkleTreeBuilder.CreateFromValues(proxy, expected, maxChildrenCount);
            var root = proofs.First().GetRoot();

            var result = await root.EnumerateItems(proxy).ToArrayAsync();

            Assert.Equal(expected, result, new ArrayComparer<Example>());
        }
    }
}