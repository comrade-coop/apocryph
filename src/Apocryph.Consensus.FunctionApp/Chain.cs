using Apocryph.HashRegistry.MerkleTree;

namespace Apocryph.Consensus
{
    public class Chain
    {
        // public Reference Creation { get; }
        public IMerkleTree<AgentState> GenesisStates { get; }
        public string ConsensusType { get; }
        public int SlotsCount { get; }

        public Chain(IMerkleTree<AgentState> genesisStates, string consensusType, int slotsCount)
        {
            GenesisStates = genesisStates;
            ConsensusType = consensusType;
            SlotsCount = slotsCount;
        }
    }
}