using Apocryph.Ipfs;

namespace Apocryph.Consensus
{
    public class AgentState
    {
        // public Reference Creation { get; private set; }
        public int Nonce { get; private set; }
        public ReferenceData Data { get; private set; }
        // public IMerkleTree<Reference> Subscriptions { get; private set; }
        public Hash<string> CodeHash { get; private set; }

        public AgentState(int nonce, ReferenceData data, Hash<string> codeHash)
        {
            Nonce = nonce;
            Data = data;
            CodeHash = codeHash;
        }
    }
}