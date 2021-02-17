namespace Apocryph.HashRegistry.MerkleTree
{
    public class MerkleTreeProof<T>
    {
        // Explanation: the Levels array stores a deep-to-shallow list of children to the right and left of the one we are prooving
        // So, for a tree which looks like this: (root has two children, each of which has 3)
        //       ____R____
        //    __1__     __2__
        //   3  V  4   5  6  7
        // The proof will look like this:
        //  Value = V
        //  Levels[0] = { Hash(3) }, { Hash(4) }
        //  Levels[1] = {}, { Hash(2) }
        // And the verification will look like this:
        //  ValueHash = Hash(V)
        //  Level0Hash = Hash(Node(Levels[0].Left + { ValueHash } + Levels[0].Right)) =
        //             = Hash(Node({Hash(3), Hash(V), Hash(5)}) = Hash(1)
        //  Level1Hash = Hash(Node(Levels[1].Left + { Level0Hash } + Levels[1].Right))
        //             = Hash(Node({Hash(1), Hash(2)}) = Hash(R)
        //  RootHash = Level1Hash

        public MerkleTreeProof(T value, (Hash<IMerkleTree<T>>[], Hash<IMerkleTree<T>>[])[] levels)
        {
            Value = value;
            Levels = levels;
        }

        public T Value { get; }

        public (Hash<IMerkleTree<T>>[] left, Hash<IMerkleTree<T>>[] right)[] Levels { get; }

        public Hash<IMerkleTree<T>> ComputeRootHash()
        {
            var previousLevelHash = Hash.From<IMerkleTree<T>>(new MerkleTreeLeaf<T>(Value));

            foreach (var (leftChildren, rightChildren) in Levels)
            {
                var children = new Hash<IMerkleTree<T>>[leftChildren.Length + 1 + rightChildren.Length];
                leftChildren.CopyTo(children, 0);
                children[leftChildren.Length] = previousLevelHash;
                rightChildren.CopyTo(children, leftChildren.Length + 1);

                previousLevelHash = Hash.From<IMerkleTree<T>>(new MerkleTreeNode<T>(children));
            }

            return previousLevelHash;
        }
    }
}