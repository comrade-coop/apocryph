using Apocryph.Ipfs;

namespace Apocryph.Consensus
{
    public class Chain
    {
        // public Reference Creation { get; private set; }
        public ChainState GenesisState { get; private set; }
        public string ConsensusType { get; private set; }
        public Hash<object>? ConsensusParameters { get; private set; }
        public int SlotsCount { get; private set; }

        public Chain(ChainState genesisState, string consensusType, Hash<object>? consensusParameters, int slotsCount)
        {
            GenesisState = genesisState;
            ConsensusType = consensusType;
            ConsensusParameters = consensusParameters;
            SlotsCount = slotsCount;
        }
    }
}