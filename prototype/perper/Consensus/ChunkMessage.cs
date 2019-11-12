using System.Collections.Generic;
using Apocryph.Execution;

namespace Apocryph.Consensus
{
    public class ChunkMessage : ConsensusMessage
    {
        // public Hash<Block> Block;
        public int Subheight;
        public Hash<ChunkMessage> PreviousChunk;
        public List<ExecutionMessage> Messages; // List<Hash<ExecutionMessage>>?
    }
}
