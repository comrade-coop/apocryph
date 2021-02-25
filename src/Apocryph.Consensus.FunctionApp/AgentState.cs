namespace Apocryph.Consensus
{
    public class AgentState
    {
        // public Reference Creation { get; }
        public ReferenceData Data { get; }
        public string Handler { get; }
        // public Hash<AgentCode> Code { get; }
        // public IMerkleTree<Subscription> Subscriptions { get; }

        public AgentState(ReferenceData data, string handler)
        {
            Data = data;
            Handler = handler;
        }
    }
}