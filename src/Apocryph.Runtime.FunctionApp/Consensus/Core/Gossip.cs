namespace Apocryph.Runtime.FunctionApp.Consensus.Core
{
    public class Gossip<T>
    {
        public T Value { get; }
        public Node[] Signers { get; }
        public GossipVerb Verb { get; }

        public Gossip(T value, Node[] signers, GossipVerb verb)
        {
            Value = value;
            Signers = signers;
            Verb = verb;
        }
    }

    public enum GossipVerb
    {
        Confirm,
        Reject
    }
}