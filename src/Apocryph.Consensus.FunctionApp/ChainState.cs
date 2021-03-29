using Apocryph.HashRegistry.MerkleTree;

namespace Apocryph.Consensus
{
    public class ChainState
    {
        public IMerkleTree<AgentState> AgentStates { get; }
        public int NextAgentNonce { get; }

        public ChainState(IMerkleTree<AgentState> agentStates, int nextAgentNonce)
        {
            AgentStates = agentStates;
            NextAgentNonce = nextAgentNonce;
        }
    }
}