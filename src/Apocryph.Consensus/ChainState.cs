using Apocryph.Ipfs.MerkleTree;

namespace Apocryph.Consensus
{
    public class ChainState
    {
        public IMerkleTree<AgentState> AgentStates { get; private set; }
        public int NextAgentNonce { get; private set; }

        public ChainState(IMerkleTree<AgentState> agentStates, int nextAgentNonce)
        {
            AgentStates = agentStates;
            NextAgentNonce = nextAgentNonce;
        }
    }
}