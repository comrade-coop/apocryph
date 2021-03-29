using Apocryph.HashRegistry;
using Apocryph.HashRegistry.MerkleTree;

namespace Apocryph.Consensus.Snowball.FunctionApp
{
    public class Block
    {
        public Hash<Block>? Previous { get; }
        public IMerkleTree<Message> InputMessages { get; }
        public IMerkleTree<Message> OutputMessages { get; }
        public ChainState State { get; }

        public Block(Hash<Block>? previous, IMerkleTree<Message> inputMessages, IMerkleTree<Message> outputMessages, ChainState state)
        {
            Previous = previous;
            InputMessages = inputMessages;
            OutputMessages = outputMessages;
            State = state;
        }
    }
}