namespace Apocryph.Consensus
{
    public class ConsensusState
    {
        public ValidatorSet ValidatorSet;
        public Hash<ChunkMessage> LastChunk;
        public int Height;
    }
}
