using Apocryph.HashRegistry.MerkleTree;

namespace Apocryph.Consensus
{
    public class Chain
    {
        // public Reference Creation { get; }
        public IMerkleTree<AgentState> GenesisStates { get; }
        public string ConsensusType { get; }

        public Chain(IMerkleTree<AgentState> genesis, string consensusType)
        {
            GenesisStates = genesis;
            ConsensusType = consensusType;
        }
    }
}