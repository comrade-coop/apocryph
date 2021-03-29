namespace Apocryph.Consensus
{
    public class Chain
    {
        // public Reference Creation { get; }
        public ChainState GenesisState { get; }
        public string ConsensusType { get; }
        public object? ConsensusParameters { get; }
        public int SlotsCount { get; }

        public Chain(ChainState genesisState, string consensusType, object? consensusParameters, int slotsCount)
        {
            GenesisState = genesisState;
            ConsensusType = consensusType;
            ConsensusParameters = consensusParameters;
            SlotsCount = slotsCount;
        }
    }
}