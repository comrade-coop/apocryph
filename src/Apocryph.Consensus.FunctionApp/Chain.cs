using Apocryph.HashRegistry.MerkleTree;

namespace Apocryph.Consensus
{
    public class Chain
    {
        // public Reference Creation { get; }
        public IMerkleTree<AgentState> GenesisStates { get; }

        public Chain(IMerkleTree<AgentState> genesis)
        {
            GenesisStates = genesis;
        }
    }
}