using Apocryph.HashRegistry;

namespace Apocryph.Consensus
{
    public class Reference
    {
        public Hash<Chain> Chain { get; }
        public Hash<AgentState> Agent { get; }
        public string[] AllowedMessageTypes { get; }
        // public MerkleTreeProof<Message> Source { get; }

        public Reference(Hash<Chain> chain, Hash<AgentState> agent, string[] allowedMessageTypes)
        {
            Chain = chain;
            Agent = agent;
            AllowedMessageTypes = allowedMessageTypes;
        }
    }
}