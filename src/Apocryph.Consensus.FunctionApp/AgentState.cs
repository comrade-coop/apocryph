using Apocryph.Ipfs;

namespace Apocryph.Consensus
{
    public class AgentState
    {
        public int Nonce { get; }
        public ReferenceData Data { get; }
        // public IMerkleTree<Subscription> Subscriptions { get; }
        public Hash<string> CodeHash { get; }

        public AgentState(int nonce, ReferenceData data, Hash<string> codeHash)
        {
            Nonce = nonce;
            Data = data;
            CodeHash = codeHash;
        }
    }
}