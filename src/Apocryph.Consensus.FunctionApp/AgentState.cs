namespace Apocryph.Consensus
{
    public class AgentState
    {
        public int Nonce { get; }
        public ReferenceData Data { get; }
        // public IMerkleTree<Subscription> Subscriptions { get; }
        public string Handler { get; }
        // public Hash<AgentCode> Code { get; }

        public AgentState(int nonce, ReferenceData data, string handler)
        {
            Nonce = nonce;
            Data = data;
            Handler = handler;
        }
    }
}