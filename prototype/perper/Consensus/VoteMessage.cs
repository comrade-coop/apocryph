namespace Apocryph.Consensus
{
    public class VoteMessage : ConsensusMessage
    {
        public int Height;
        public Hash<ChunkMessage> LastValid;
        // public Block ForBlock;
    }
}
