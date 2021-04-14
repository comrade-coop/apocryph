using Apocryph.Ipfs;
using Apocryph.Ipfs.MerkleTree;

namespace Apocryph.Consensus.Snowball.FunctionApp
{
    public class Block
    {
        // public int Height;
        public Hash<Block>? Previous { get; private set; }
        public IMerkleTree<Message> InputMessages { get; private set; }
        public IMerkleTree<Message> OutputMessages { get; private set; }
        public ChainState State { get; private set; }

        public Block(Hash<Block>? previous, IMerkleTree<Message> inputMessages, IMerkleTree<Message> outputMessages, ChainState state)
        {
            Previous = previous;
            InputMessages = inputMessages;
            OutputMessages = outputMessages;
            State = state;
        }
    }
}